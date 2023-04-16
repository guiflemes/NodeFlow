package adapters

import (
	"context"
	"database/sql"
	"encoding/json"
	"flowChart/domain"
	"fmt"
	"log"
	"sync"
	"time"

	"flowChart/settings"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DatabaseConfig struct {
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

func (conf *DatabaseConfig) Parse() {
	conf.User = settings.GETENV("POSTGRES_USER")
	conf.Password = settings.GETENV("POSTGRES_PASSWORD")
	conf.Host = settings.GETENV("POSTGRES_HOST")
	conf.Port = settings.GETENV("POSTGRES_PORT")
	conf.Database = settings.GETENV("POSTGRES_DB_NAME")
}

type BaseFlowChartAggregate[T comparable] struct {
	client *sqlx.DB
}

func NewBaseFlowchartAggregate[T comparable](conf *DatabaseConfig) *BaseFlowChartAggregate[T] {
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

	return &BaseFlowChartAggregate[T]{
		client: mydb,
	}
}

func ToJsonB(value any) []byte {
	result, _ := json.Marshal(value)
	return result
}

func (r *BaseFlowChartAggregate[T]) StoreFlowChart(ctx context.Context, flowChart *domain.FlowChart[T]) error {
	query := `INSERT into flowchart (title, key) VALUES ($1, $2) RETURNING id`
	stmt, err := r.client.PrepareContext(ctx, query)

	if err != nil {
		return fmt.Errorf("error to prepare flowchart stmt: %w", err)
	}

	err = stmt.QueryRowContext(ctx, flowChart.Title, flowChart.Key).Scan(&flowChart.Id)

	if err != nil {
		return fmt.Errorf("error storing an user: %w", err)
	}

	return r.createOrUpdateNode(ctx, flowChart)

}

func (r *BaseFlowChartAggregate[T]) UpdateFlowChart(ctx context.Context, flowChart *domain.FlowChart[T]) error {
	query := "UPDATE flowchart set title=$2 WHERE key=$1"

	stmt, err := r.client.PrepareContext(ctx, query)

	if err != nil {
		return fmt.Errorf("error preparing stmt to update a flowchart: %w", err)
	}

	result, err := stmt.ExecContext(ctx, flowChart.Title, flowChart.Key)

	if err != nil {
		return fmt.Errorf("error updating a flowchart: %w", err)
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		return fmt.Errorf("there is no flowchart to the given key")
	}

	return r.createOrUpdateNode(ctx, flowChart)
}

func (r *BaseFlowChartAggregate[T]) FlowChartExists(ctx context.Context, flowChartKey string) (bool, error) {
	var exists bool
	query := "SELECT * FROM flowchart WHERE key=$1"

	err := r.client.QueryRowContext(ctx, query, flowChartKey).Scan(&exists)

	if err == nil {
		return exists, nil
	}

	if err.Error() == "sql: no rows in result set" {
		return false, nil
	}

	return false, err

}

func (r *BaseFlowChartAggregate[T]) undoNode(ctx context.Context, flowchartID string, internal_id string) error {
	query := "DELETE FROM node WHERE flowchart=$1 AND internal_id=$2"

	stmt, err := r.client.PrepareContext(ctx, query)

	if err != nil {
		return fmt.Errorf("error preparing stmt to delete a node: %w", err)
	}

	if _, err := stmt.ExecContext(ctx, flowchartID, internal_id); err != nil {
		return fmt.Errorf("error deleting a node: %w", err)
	}

	return nil
}

func (r *BaseFlowChartAggregate[T]) createOrUpdateNode(ctx context.Context, flowChart *domain.FlowChart[T]) error {

	saveFunc := func(n *domain.Node[T]) error {
		exists, err := r.nodeExists(ctx, flowChart.Id, n)

		if err != nil {
			return err
		}

		if exists {
			return nil
		}

		return r.createNode(ctx, flowChart.Id, n)
	}

	undo := make([]*domain.Node[T], 0)

	flowChart.Node.Traverse(domain.TraversePreOrder, domain.TraverseAll, -1, func(n *domain.Node[T]) bool {
		err := saveFunc(n)

		if err != nil { // Stop traversing the tree if an error occurs, add parent to undo
			undo = append(undo, n)
			return true
		}

		return false
	})

	if len(undo) > 0 {
		var wg sync.WaitGroup
		for _, n := range undo {
			wg.Add(1)
			go func(nodeID string) {
				defer wg.Done()
				r.undoNode(ctx, flowChart.Id, nodeID)
			}(n.NodeID)
		}
		wg.Wait()
	}

	return nil

}

func (r *BaseFlowChartAggregate[T]) nodeExists(ctx context.Context, flowchartID string, node *domain.Node[T]) (bool, error) {
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

func (r *BaseFlowChartAggregate[T]) updateNode(ctx context.Context, flowchartID string, node *domain.Node[T]) error {
	query := `
	UPDATE node set 
	parent_id=$3,
	dragging=$4,
	selected=$5,
	position_absolute=$6,
	height=$7,
	width=$8,
	position=$9,
	data=$10
	WHERE id=$1 AND internal_id=$2`

	stmt, err := r.client.PrepareContext(ctx, query)

	if err != nil {
		return fmt.Errorf("error preparing stmt to update a node: %w", err)
	}

	result, err := stmt.ExecContext(ctx,
		flowchartID,
		node.NodeID,
		node.ParentId(),
		node.Dragging,
		node.Selected,
		ToJsonB(node.PositionAbsolute),
		node.Height,
		node.Width,
		ToJsonB(node.Position),
		ToJsonB(node.Data),
	)

	if err != nil {
		return fmt.Errorf("error updating a node: %w", err)
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		return fmt.Errorf(`there is no node to the given flowchartID "%s" and nodeID "%s": %w`, flowchartID, node.NodeID, err)
	}

	return nil
}

func (r *BaseFlowChartAggregate[T]) createNode(ctx context.Context, flowchartID string, node *domain.Node[T]) error {
	query := `INSERT into node (internal_id, parent_id, flowchart_id, dragging, selected, position_absolute, height, width, position, data)
	 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`

	stmt, err := r.client.PrepareContext(ctx, query)

	if err != nil {
		return fmt.Errorf("error preparing stmt to Store a node: %w", err)
	}

	defer stmt.Close()

	if _, err := stmt.ExecContext(
		ctx,
		node.NodeID,
		node.ParentId(),
		flowchartID,
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
