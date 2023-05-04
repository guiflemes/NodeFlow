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

	flowChartUnstructuredDataAgr := adapters.NewFlowChartUnstructuredDataAgg(newPsqlClient)

	editFlowChart := command.NewHandlerFlowChartUnstructuredData(flowChartUnstructuredDataAgr)

	return handlers.Application{
		Commands: handlers.Commands{
			EditFlowChart: editFlowChart,
		},
	}
}
