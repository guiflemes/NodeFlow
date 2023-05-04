package handlers

import "flowChart/handlers/command"

type Commands struct {
	EditFlowChart command.HandlerFlowChartUnstructuredData
}

type Application struct {
	Commands Commands
}
