package service

import (
	"flowChart/adapters"
	"flowChart/handlers"
	"flowChart/handlers/command"
)

func Bootstrap() handlers.Application {
	config := &DatabaseConfig{}
	config.Parse()
	newPsqlClient := NewPostgresDb(config)

	flowChartSimpleDataAgr := adapters.NewFlowChartDataRepo(newPsqlClient)

	editFlowChart := command.NewHandlerFlowChartSimpleData(flowChartSimpleDataAgr)

	return handlers.Application{
		Commands: handlers.Commands{
			EditFlowChart: editFlowChart,
		},
	}
}
