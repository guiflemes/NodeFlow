package transport

type DataDto struct {
	Label string `json:"label"`
}

type UnstructuredDataDto interface{}

type FlowChartDto[T comparable] struct {
	Title string        `json:"title"`
	Key   string        `json:"key"`
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
	Selected         bool        `json:"selected"`
	PositionAbsolute PositionDto `json:"positionAbsolute"`
	Dragging         bool        `json:"dragging"`
	Type             string      `json:"type"`
}

type EdgeDto struct {
	Id     string `json:"id"`
	Source string `json:"source"`
	Target string `json:"target"`
}

var FlowChartJson string = `{
	"title": "First Flow",
	"key" : "first_flow",
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
