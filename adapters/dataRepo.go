package adapters

import "flowChart/ports"

type FlowChartDataAggregate struct {
	*BaseFlowChartAggregate[ports.Data]
}

func NewFlowChartDataRepo(config *DatabaseConfig) *FlowChartDataAggregate {
	return &FlowChartDataAggregate{
		BaseFlowChartAggregate: NewBaseFlowchartAggregate[ports.Data](config),
	}
}
