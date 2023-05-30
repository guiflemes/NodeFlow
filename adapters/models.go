package adapters

import (
	"encoding/json"
	"errors"
)

type PositionModel struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

func (p *PositionModel) Scan(value any) error {
	b, ok := value.([]uint8)
	if !ok {
		return errors.New("type assertion to []uint8 failed")
	}

	return json.Unmarshal(b, &p)
}

type EdgeModel struct {
	Id     int    `json:"id"`
	Source string `json:"source"`
	Target string `json:"target"`
}

type NodeModel[T any] struct {
	NodeID           string        `json:"id" db:"internal_id"`
	ParentID         string        `json:"parentID" db:"parent_id"`
	Position         PositionModel `json:"position" db:"position"`
	Data             T             `json:"data" db:"data"`
	Width            int16         `json:"width" db:"width"`
	Height           int16         `json:"height" db:"height"`
	Selected         bool          `json:"selected" db:"selected"`
	PositionAbsolute PositionModel `json:"position_absolute" db:"position_absolute"`
	Dragging         bool          `json:"dragging" db:"dragging"`
	Type             string        `json:"type" db:"type"`
}

type FlowChartModel[T any] struct {
	ID    string          `json:"id" db:"id"`
	Title string          `json:"title" db:"title"`
	Key   string          `json:"key" db:"key"`
	Nodes []*NodeModel[T] `json:"nodes"`
	Edges []*EdgeModel    `json:"Edges"`
}

func (f *FlowChartModel[T]) AddNode(node *NodeModel[T]) {
	f.Nodes = append(f.Nodes, node)
	edge := f.createEdge(node)
	if edge != nil {
		f.Edges = append(f.Edges, edge)
	}

}

func (f *FlowChartModel[T]) createEdge(node *NodeModel[T]) *EdgeModel {

	if node.NodeID == node.ParentID {
		return nil
	}

	return &EdgeModel{
		Id:     len(f.Edges) + 1,
		Source: node.NodeID,
		Target: node.ParentID,
	}

}

type WagtailDataModel []map[string]any

func (u *WagtailDataModel) Scan(value any) error {
	b, ok := value.([]uint8)
	if !ok {
		return errors.New("type assertion to []uint8 failed")
	}

	return json.Unmarshal(b, &u)
}
