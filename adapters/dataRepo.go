package adapters

import (
	"flowChart/domain"
)

type FlowChartDataAggregate struct {
	*BaseFlowChartAggregate[domain.Data]
}

func NewFlowChartDataRepo(config *DatabaseConfig) *FlowChartDataAggregate {
	return &FlowChartDataAggregate{
		BaseFlowChartAggregate: NewBaseFlowchartAggregate[domain.Data](config),
	}
}
