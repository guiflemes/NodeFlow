package adapters

import (
	"context"
	"encoding/json"
	"errors"
	"flowChart/domain"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type BaseFlowChartAggregate[T any] struct {
	client *sqlx.DB
}

func NewBaseFlowchartAggregate[T any](db *sqlx.DB) *BaseFlowChartAggregate[T] {
	return &BaseFlowChartAggregate[T]{
		client: db,
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

	return r.editNode(ctx, flowChart)

}

func (r *BaseFlowChartAggregate[T]) UpdateFlowChart(ctx context.Context, flowChart *domain.FlowChart[T]) error {
	query := `UPDATE flowchart set title=$1 WHERE key=$2`

	stmt, err := r.client.PrepareContext(ctx, query)

	if err != nil {
		return fmt.Errorf("error preparing stmt to update a flowchart: %w", err)
	}

	result, err := stmt.ExecContext(ctx, flowChart.Title, flowChart.Key)

	if err != nil {
		return fmt.Errorf("error updating a flowchart: %w", err)
	}

	if rows, err := result.RowsAffected(); rows == 0 {
		return fmt.Errorf("there is no flowchart to the given key: %w", err)
	}

	return r.editNode(ctx, flowChart)
}

func (r *BaseFlowChartAggregate[T]) FlowChartExists(ctx context.Context, flowChart *domain.FlowChart[T]) (bool, error) {
	query := "SELECT id FROM flowchart WHERE key=$1"

	err := r.client.QueryRowContext(ctx, query, flowChart.Key).Scan(&flowChart.Id)

	if err == nil {
		return flowChart.Id != "", nil
	}

	if err.Error() == "sql: no rows in result set" {
		return false, nil
	}

	return false, fmt.Errorf("error checking if flowchart exists: %w", err)

}

func (r *BaseFlowChartAggregate[T]) deleteNodes(ctx context.Context, tx *sqlx.Tx, flowchartID string) error {
	query := "DELETE FROM node WHERE flowchart_id=$1"

	stmt, err := tx.PrepareContext(ctx, query)

	if err != nil {
		return fmt.Errorf("error preparing stmt to delete a node: %w", err)
	}

	if _, err := stmt.ExecContext(ctx, flowchartID); err != nil {
		return fmt.Errorf("error deleting a node: %w", err)
	}

	return nil
}

func (r *BaseFlowChartAggregate[T]) createNode(ctx context.Context, tx *sqlx.Tx, flowchartID string, node *domain.Node[T]) error {
	query := `INSERT into node (internal_id, parent_id, flowchart_id, dragging, selected, position_absolute, height, width, position, data, type)
	 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id`

	stmt, err := tx.PrepareContext(ctx, query)

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
		node.PositionAbsolute,
		node.Height,
		node.Width,
		node.Position,
		ToJsonB(node.Data),
		node.Type,
	); err != nil {
		return fmt.Errorf("error creating a node: %w", err)
	}

	return nil
}

func (r *BaseFlowChartAggregate[T]) editNode(ctx context.Context, flowChart *domain.FlowChart[T]) error {

	saveFunc := func(tx *sqlx.Tx, n *domain.Node[T]) error {
		return r.createNode(ctx, tx, flowChart.Id, n)
	}

	var errR error
	handlerErr := func(err error) {
		if err == nil {
			return
		}

		if errR != nil {
			errR = errors.Join(errR, err)
			return
		}

		errR = err
	}

	return r.RunInTransaction(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		if err := r.deleteNodes(ctx, tx, flowChart.Id); err != nil {
			return err
		}

		flowChart.Node.Traverse(domain.TraversePreOrder, domain.TraverseAll, -1, func(n *domain.Node[T]) bool {
			err := saveFunc(tx, n)
			handlerErr(err)
			return false
		})

		return errR

	})
}

func (r *BaseFlowChartAggregate[T]) RunInTransaction(ctx context.Context, txFunc func(ctx context.Context, tx *sqlx.Tx) error) error {
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

	err = txFunc(ctx, tx)

	if err != nil {
		return err
	}

	return nil
}

func (r *BaseFlowChartAggregate[T]) GetFlowChart(ctx context.Context, key string) (*FlowChartModel[T], error) {
	query := `
	SELECT
		flow.id,
		flow.title,
		flow.key,
		node.internal_id,
		node.parent_id,
		node.position,
  		node.data,
		node.width,
		node.height,
		node.position_absolute,
		node.selected,
		node.dragging,
		node.type
	FROM
		node
	JOIN
		(
		SELECT *
		FROM flowchart
		WHERE flowchart.key = $1
	) as flow
	ON
	flow.id = node.flowchart_id
	`

	flow := &FlowChartModel[T]{}

	rows, err := r.client.QueryxContext(ctx, query, key)

	if err != nil {
		return flow, fmt.Errorf("error querying a flowchart: %w", err)
	}

	for rows.Next() {
		node := &NodeModel[T]{}
		var (
			flowchartID    string
			flowchartTitle string
			flowchartKey   string
		)
		if err := rows.Err(); err != nil {
			return flow, fmt.Errorf("error querying a flowchart %w", err)
		}

		err := rows.Scan(
			&flowchartID,
			&flowchartTitle,
			&flowchartKey,
			&node.NodeID,
			&node.ParentID,
			&node.Position,
			&node.Data,
			&node.Width,
			&node.Height,
			&node.PositionAbsolute,
			&node.Selected,
			&node.Dragging,
			&node.Type,
		)
		if err != nil {
			return flow, fmt.Errorf("error querying a flowchart %w", err)
		}

		if flow.ID == "" {
			flow.ID = flowchartID
			flow.Title = flowchartTitle
			flow.Key = flowchartKey
		}

		flow.AddNode(node)

	}

	return flow, nil

}
