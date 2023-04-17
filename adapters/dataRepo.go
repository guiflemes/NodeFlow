package adapters

import (
	"flowChart/domain"

	"github.com/jmoiron/sqlx"
)

type FlowChartDataAggregate struct {
	*BaseFlowChartAggregate[domain.Data]
}

func NewFlowChartDataRepo(client *sqlx.DB) *FlowChartDataAggregate {
	return &FlowChartDataAggregate{
		BaseFlowChartAggregate: NewBaseFlowchartAggregate[domain.Data](client),
	}
}
