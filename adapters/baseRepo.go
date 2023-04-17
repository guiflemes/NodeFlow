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

type NodeError []error

func (e *NodeError) Error() string {
	return ""
}

type BaseFlowChartAggregate[T comparable] struct {
	client *sqlx.DB
}

func NewBaseFlowchartAggregate[T comparable](db *sqlx.DB) *BaseFlowChartAggregate[T] {
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
	query := `INSERT into node (internal_id, parent_id, flowchart_id, dragging, selected, position_absolute, height, width, position, data)
	 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`

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

// func (r *BaseFlowChartAggregate[T]) undoNode(ctx context.Context, flowchartID string, internal_id string) error {
// 	query := "DELETE FROM node WHERE flowchart_id=$1 AND internal_id=$2"

// 	stmt, err := r.client.PrepareContext(ctx, query)

// 	if err != nil {
// 		return fmt.Errorf("error preparing stmt to delete a node: %w", err)
// 	}

// 	if _, err := stmt.ExecContext(ctx, flowchartID, internal_id); err != nil {
// 		return fmt.Errorf("error deleting a node: %w", err)
// 	}

// 	return nil
// }

// func (r *BaseFlowChartAggregate[T]) createOrUpdateNode(ctx context.Context, flowChart *domain.FlowChart[T]) error {

// 	saveFunc := func(n *domain.Node[T]) *NodeError {

// 		exists, err := r.nodeExists(ctx, flowChart.Id, n)

// 		if err != nil {
// 			return err
// 		}

// 		if exists {
// 			return r.updateNode(ctx, flowChart.Id, n)
// 		}

// 		return r.createNode(ctx, flowChart.Id, n)
// 	}

// 	undo := make([]*domain.Node[T], 0)

// 	flowChart.Node.Traverse(domain.TraversePreOrder, domain.TraverseAll, -1, func(n *domain.Node[T]) bool {
// 		err := saveFunc(n)

// 		if err != nil {
// 			logrus.Info(fmt.Sprintf("Stop traversing the tree to node %s, err: %s", n.NodeID, err))

// 			if err.code == CreateErr {
// 				logrus.Info(fmt.Sprintf("Error on node creating, undo node  %s and its children", n.NodeID))
// 				undo = append(undo, n)
// 			}
// 			return true
// 		}

// 		return false
// 	})

// 	if len(undo) > 0 {
// 		var wg sync.WaitGroup
// 		for _, n := range undo {
// 			wg.Add(1)
// 			go func(nodeID string) {
// 				defer wg.Done()
// 				r.undoNode(ctx, flowChart.Id, nodeID)
// 			}(n.NodeID)
// 		}
// 		wg.Wait()
// 	}

// 	return nil

// }

// func (r *BaseFlowChartAggregate[T]) nodeExists(ctx context.Context, flowchartID string, node *domain.Node[T]) (bool, *NodeError) {
// 	var id string
// 	query := "SELECT id FROM node WHERE flowchart_id=$1 AND internal_id=$2"

// 	err := r.client.QueryRowContext(ctx, query, flowchartID, node.NodeID).Scan(&id)

// 	if err == nil {
// 		return id != "", nil
// 	}

// 	if err.Error() == "sql: no rows in result set" {
// 		return false, nil
// 	}

// 	logrus.WithError(err)
// 	return false, &NodeError{err: fmt.Errorf("error checking if node exists: %w", err), code: GetErr}
// }

// func (r *BaseFlowChartAggregate[T]) updateNode(ctx context.Context, flowchartID string, node *domain.Node[T]) *NodeError {
// 	// TODO VER PQ NAO ATUALIZA
// 	query := `
// 	UPDATE node set
// 	parent_id=$3,
// 	dragging=$4,
// 	selected=$5,
// 	position_absolute=$6,
// 	height=$7,
// 	width=$8,
// 	position=$9,
// 	data=$10
// 	WHERE id=$1 AND internal_id=$2`

// 	stmt, err := r.client.PrepareContext(ctx, query)

// 	if err != nil {
// 		return &NodeError{err: fmt.Errorf("error preparing stmt to update a node: %w", err), code: UpdateErr}
// 	}

// 	result, err := stmt.ExecContext(ctx,
// 		flowchartID,
// 		node.NodeID,
// 		node.ParentId(),
// 		node.Dragging,
// 		node.Selected,
// 		ToJsonB(node.PositionAbsolute),
// 		node.Height,
// 		node.Width,
// 		ToJsonB(node.Position),
// 		ToJsonB(node.Data),
// 	)

// 	if err != nil {
// 		return &NodeError{err: fmt.Errorf("error updating a node: %w", err), code: UpdateErr}
// 	}

// 	if rows, err := result.RowsAffected(); rows == 0 {
// 		return &NodeError{err: fmt.Errorf(`there is no node to the given flowchartID "%s" and nodeID "%s": %w`, flowchartID, node.NodeID, err), code: UpdateErr}
// 	}

// 	return nil
// }
