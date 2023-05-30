package adapters

import (
	"flowChart/domain"

	"github.com/jmoiron/sqlx"
)

type WriteFlowChartUnstructuredDataAgg struct {
	*BaseFlowChartAggregate[domain.UnstructuredDataDomain]
}

func NewWriteFlowChartUnstructuredDataAgg(client *sqlx.DB) *WriteFlowChartUnstructuredDataAgg {
	return &WriteFlowChartUnstructuredDataAgg{
		BaseFlowChartAggregate: NewBaseFlowchartAggregate[domain.UnstructuredDataDomain](client),
	}
}

type ReadFlowChartUnstructuredDataAgg struct {
	*BaseFlowChartAggregate[WagtailDataModel]
}

func NewReadFlowChartUnstructuredDataAgg(client *sqlx.DB) *ReadFlowChartUnstructuredDataAgg {
	return &ReadFlowChartUnstructuredDataAgg{
		BaseFlowChartAggregate: NewBaseFlowchartAggregate[WagtailDataModel](client),
	}
}
