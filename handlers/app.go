package handlers

import "flowChart/handlers/command"

type Commands struct {
	EditFlowChart command.HandlerFlowChartSimpleData
}

type Application struct {
	Commands Commands
}
