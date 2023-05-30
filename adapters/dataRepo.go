package adapters

import (
	"flowChart/domain"

	"github.com/jmoiron/sqlx"
)

// use as example
type FlowChartDataAggregate struct {
	*BaseFlowChartAggregate[domain.Data]
}

func NewFlowChartDataRepo(client *sqlx.DB) *FlowChartDataAggregate {
	return &FlowChartDataAggregate{
		BaseFlowChartAggregate: NewBaseFlowchartAggregate[domain.Data](client),
	}
}
