package main

import (
	"errors"
	"flowChart/domain"
	"sync"
)

type FlowChartDto[T comparable] struct {
	Title string        `json:"title"`
	Nodes []*NodeDto[T] `json:"nodes"`
	Edges []*EdgeDto    `json:"Edges"`
}

type PositionDto struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type NodeDto[T comparable] struct {
	Id               string      `json:"id"`
	Position         PositionDto `json:"position"`
	Data             T           `json:"data"`
	Width            int16       `json:"width"`
	Height           int16       `json:"height"`
	Selected         bool        `json:"selected "`
	PositionAbsolute PositionDto `json:"positionAbsolute"`
	Dragging         bool        `json:"dragging"`
}

type EdgeDto struct {
	Id     string `json:"id"`
	Source string `json:"source"`
	Target string `json:"target"`
}

func toDomain[T comparable](flowChart *FlowChartDto[T]) (*domain.FlowChart[T], error) {
	m := sync.RWMutex{}

	nodeMap := func() map[string]*domain.Node[T] {
		nodeM := make(map[string]*domain.Node[T], len(flowChart.Nodes))
		for _, n := range flowChart.Nodes {
			m.Lock()
			nodeM[n.Id] = domain.NewNode(n.Id, n.Data, domain.Position{X: n.Position.X, Y: n.Position.Y},
				n.Width, n.Height, n.Selected, domain.Position{X: n.PositionAbsolute.X, Y: n.PositionAbsolute.Y}, n.Dragging)
			m.Unlock()
		}
		return nodeM
	}()

	var flow *domain.FlowChart[T]
	var root *domain.Node[T]

	for _, edge := range flowChart.Edges {

		if err := func() error {
			m.RLock()
			defer m.RUnlock()
			parent, ok := nodeMap[edge.Source]

			if !ok {
				return errors.New("Parent Not Found")
			}

			if edge.Source == "0" {
				root = parent
			}

			if child, ok := nodeMap[edge.Id]; ok {
				parent.AddChild(child)
				return nil
			}

			return errors.New("Child Not Found")

		}(); err != nil {
			return flow, err
		}

	}

	return &domain.FlowChart[T]{Title: flowChart.Title, Node: root}, nil
}

var flowChartJson string = `{
	"title": "firstFlow",
	"nodes": [
	{
		"id": "0",
		"type": "input",
		"data": {
			"label": "Node"
		},
		"position": {
			"x": -55,
			"y": -68.5
		},
		"width": 150,
		"height": 40,
		"selected": false,
		"positionAbsolute": {
			"x": -55,
			"y": -68.5
		},
		"dragging": false
	},
	{
		"id": "1",
		"position": {
			"x": -195.25,
			"y": 21.759944915771484
		},
		"data": {
			"label": "Node 1"
		},
		"width": 150,
		"height": 40,
		"selected": false,
		"positionAbsolute": {
			"x": -195.25,
			"y": 21.759944915771484
		},
		"dragging": false
	},
	{
		"id": "2",
		"position": {
			"x": 156.75,
			"y": 16.759944915771484
		},
		"data": {
			"label": "Node 2"
		},
		"width": 150,
		"height": 40,
		"selected": false,
		"positionAbsolute": {
			"x": 156.75,
			"y": 16.759944915771484
		},
		"dragging": false
	},
	{
		"id": "3",
		"position": {
			"x": -292.75,
			"y": 105.25994491577148
		},
		"data": {
			"label": "Node 3"
		},
		"width": 150,
		"height": 40,
		"selected": false,
		"positionAbsolute": {
			"x": -292.75,
			"y": 105.25994491577148
		},
		"dragging": false
	},
	{
		"id": "4",
		"position": {
			"x": -110.25,
			"y": 103.25994491577148
		},
		"data": {
			"label": "Node 4"
		},
		"width": 150,
		"height": 40,
		"selected": true,
		"positionAbsolute": {
			"x": -110.25,
			"y": 103.25994491577148
		},
		"dragging": false
	},
	{
		"id": "5",
		"position": {
			"x": 101.25,
			"y": 100.75994491577148
		},
		"data": {
			"label": "Node 5"
		},
		"width": 150,
		"height": 40,
		"selected": false,
		"positionAbsolute": {
			"x": 101.25,
			"y": 100.75994491577148
		},
		"dragging": false
	},
	{
		"id": "6",
		"position": {
			"x": 334.75,
			"y": 114.25994491577148
		},
		"data": {
			"label": "Node 6"
		},
		"width": 150,
		"height": 40
	}
],
"edges": [
	{
		"id": "1",
		"source": "0",
		"target": "1"
	},
	{
		"id": "2",
		"source": "0",
		"target": "2"
	},
	{
		"id": "3",
		"source": "1",
		"target": "3"
	},
	{
		"id": "4",
		"source": "1",
		"target": "4"
	},
	{
		"id": "5",
		"source": "2",
		"target": "5"
	},
	{
		"id": "6",
		"source": "2",
		"target": "6"
	}
]
}`
