package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"flowChart/domain"
	"fmt"
	"log"
	"time"

	"flowChart/settings"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Database struct {
	User               string
	Password           string
	Host               string
	Port               string
	Database           string
	IdleConnection     int
	OpenConnection     int
	ConnectionLifeTime int
	ConnectionIdleTime int
	ReadTimeout        int
	WriteTimeout       int
	Timeout            int
}

func (db *Database) Parse() {
	db.User = settings.GETENV("POSTGRES_USER")
	db.Password = settings.GETENV("POSTGRES_PASSWORD")
	db.Host = settings.GETENV("POSTGRES_HOST")
	db.Port = settings.GETENV("POSTGRES_PORT")
	db.Database = settings.GETENV("POSTGRES_DB_NAME")
}

type txFlowchart[T comparable] struct {
	transaction *sqlx.Tx
	flowChartID string
	node        *domain.Node[T]
}

type PostgresRepo[T comparable] struct {
	client *sqlx.DB
}

func NewPostgresRepo[T comparable](conf *Database) *PostgresRepo[T] {
	dns := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s connect_timeout=30 sslmode=disable",
		conf.Host,
		conf.Port,
		conf.User,
		conf.Password,
		conf.Database,
	)

	db, err := sql.Open("postgres", dns)
	if err != nil {
		log.Fatal(err)
	}

	mydb := sqlx.NewDb(db, "postgres")

	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}

	mydb.SetMaxOpenConns(conf.OpenConnection)
	mydb.SetMaxIdleConns(conf.IdleConnection)
	mydb.SetConnMaxLifetime(time.Duration(30) * time.Millisecond)
	mydb.SetConnMaxIdleTime(time.Duration(30) * time.Millisecond)

	return &PostgresRepo[T]{
		client: mydb,
	}
}

func (r *PostgresRepo[T]) Store(ctx context.Context, flowChart *domain.FlowChart[T]) error {
	query := `INSERT into flowchart (title) VALUES ($1) RETURNING id`
	stmt, err := r.client.PrepareContext(ctx, query)

	if err != nil {
		return fmt.Errorf("error to prepare flowchart stmt: %w", err)
	}

	err = stmt.QueryRowContext(ctx, flowChart.Title).Scan(&flowChart.Id)

	if err != nil {
		return fmt.Errorf("error storing an user: %w", err)
	}

	err = r.RunInTransaction(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		return r.storeNode(ctx, flowChart)

	})

	return err

}

//erro na query
//chan nao recebe err?
//Rollback nao deleta o flowchart

func (r *PostgresRepo[T]) storeNode(ctx context.Context, flowChart *domain.FlowChart[T]) error {
	// errorChan := make(chan error)
	var e error

	flowChart.Node.Traverse(domain.TraversePreOrder, domain.TraverseAll, -1, func(n *domain.Node[T]) bool {

		// todo GOROUTINE RAISES LOGS POSTGRES Could not receive data from client: Connection reset by pee
		func() {
			exists, err := r.nodeExists(ctx, flowChart.Id, n)

			if err != nil {
				e = err
			}

			if exists {
				return
			}

			// 	if err := r.updateNode(ctx, tx); err != nil {
			// 		errorChan <- err
			// 		return
			// 	}

			// }
			if err := r.createNode(ctx, flowChart.Id, n); err != nil {
				e = err
				// fmt.Println("err", err)
				// errorChan <- err
			}

			fmt.Println("Pos Create Node")

		}()

		return false
	})

	return e

	// select {
	// case err := <-errorChan:
	// 	return err
	// default:
	// 	return nil
	// }
}

func (r *PostgresRepo[T]) nodeExists(ctx context.Context, flowchartID string, node *domain.Node[T]) (bool, error) {
	var exists bool
	query := "SELECT * FROM node WHERE id=$1 AND internal_id=$2"

	err := r.client.QueryRowContext(ctx, query, flowchartID, node.NodeID).Scan(&exists)

	if err == nil {
		return exists, nil
	}

	if err.Error() == "sql: no rows in result set" {
		return false, nil
	}

	return false, err
}

func (r *PostgresRepo[T]) flowExists(ctx context.Context, flowchartID string) (bool, error) {
	var exists bool
	query := "SELECT * FROM flowchart WHERE id=$1"

	err := r.client.QueryRowContext(ctx, query, flowchartID).Scan(&exists)

	if err == nil {
		return exists, nil
	}

	if err.Error() == "sql: no rows in result set" {
		return false, nil
	}

	return false, err
}

func (r *PostgresRepo[T]) updateNode(ctx context.Context, tx *txFlowchart[T]) error {
	return nil
}

func (r *PostgresRepo[T]) createNode(ctx context.Context, flowchartID string, node *domain.Node[T]) error {
	query := `INSERT into node (internal_id, parent_id, flowchart_id, dragging, selected, position_absolute, height, width, position, data)
	 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`

	stmt, err := r.client.PrepareContext(ctx, query)

	if err != nil {
		return fmt.Errorf("error preparing stmt to Store a node: %w", err)
	}

	defer stmt.Close()

	ToJsonB := func(value any) []byte {
		result, _ := json.Marshal(value)
		return result
	}

	if _, err := stmt.ExecContext(
		ctx,
		node.NodeID,
		node.ParentId(),
		flowchartID, //"89068c6b-51d3-427e-a9a1-2c02dbe17f9a"
		node.Dragging,
		node.Selected,
		ToJsonB(node.PositionAbsolute),
		node.Height,
		node.Width,
		ToJsonB(node.Position),
		ToJsonB(node.Data),
	); err != nil {
		return fmt.Errorf("error creating a node: %w", err)
	}

	return nil
}

func (r *PostgresRepo[T]) RunInTransaction(ctx context.Context, txFunc func(ctx context.Context, tx *sqlx.Tx) error) error {
	tx, err := r.client.Beginx()

	if err != nil {
		return fmt.Errorf("error beginning a transaction: %w", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	if err := txFunc(ctx, tx); err != nil {
		return err
	}

	return nil
}
