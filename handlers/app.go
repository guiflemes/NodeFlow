package handlers

import (
	"flowChart/handlers/command"
	"flowChart/handlers/queries"
)

type Commands struct {
	EditFlowChart command.HandlerFlowChartUnstructuredData
}

type Queries struct {
	GetFlowChart queries.HandlerGetFlowChartUnstructuredData
}

type Application struct {
	Commands Commands
	Queries  Queries
}
