package service

import (
	"flowChart/adapters"
	"flowChart/handlers"
	"flowChart/handlers/command"
	"flowChart/handlers/queries"
)

func Bootstrap() handlers.Application {
	config := &DatabaseConfig{}
	config.Parse()
	newPsqlClient := NewPostgresDb(config)

	writeFlowChartUnstructuredDataAgr := adapters.NewWriteFlowChartUnstructuredDataAgg(newPsqlClient)
	readFlowChartUnstructuredDataAgr := adapters.NewReadFlowChartUnstructuredDataAgg(newPsqlClient)

	editFlowChart := command.NewHandlerFlowChartUnstructuredData(writeFlowChartUnstructuredDataAgr)
	getFlowChart := queries.NewHandlerGetFlowChartUnstructuredData(readFlowChartUnstructuredDataAgr)

	return handlers.Application{
		Commands: handlers.Commands{
			EditFlowChart: editFlowChart,
		},
		Queries: handlers.Queries{
			GetFlowChart: getFlowChart,
		},
	}
}
