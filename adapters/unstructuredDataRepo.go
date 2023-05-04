package adapters

import (
	"flowChart/domain"

	"github.com/jmoiron/sqlx"
)

type FlowChartUnstructuredDataAgg struct {
	*BaseFlowChartAggregate[domain.UnstructuredDataDomain]
}

func NewFlowChartUnstructuredDataAgg(client *sqlx.DB) *FlowChartUnstructuredDataAgg {
	return &FlowChartUnstructuredDataAgg{
		BaseFlowChartAggregate: NewBaseFlowchartAggregate[domain.UnstructuredDataDomain](client),
	}
}
