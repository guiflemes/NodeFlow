package queries

import (
	"context"
	"flowChart/adapters"
)

type QueryFlowChartAggregate[T any] interface {
	GetFlowChart(ctx context.Context, key string) (*adapters.FlowChartModel[T], error)
}

type HandlerGetFlowChart[T any] struct {
	agg QueryFlowChartAggregate[T]
}

func NewGetFlowChartHandler[T any](agg QueryFlowChartAggregate[T]) *HandlerGetFlowChart[T] {
	return &HandlerGetFlowChart[T]{
		agg: agg,
	}
}

func (h *HandlerGetFlowChart[T]) Handler(ctx context.Context, key string) (*adapters.FlowChartModel[T], error) {
	return h.agg.GetFlowChart(ctx, key)
}

type HandlerGetFlowChartUnstructuredData struct {
	*HandlerGetFlowChart[adapters.WagtailDataModel]
}

func NewHandlerGetFlowChartUnstructuredData(agr *adapters.ReadFlowChartUnstructuredDataAgg) HandlerGetFlowChartUnstructuredData {
	return HandlerGetFlowChartUnstructuredData{
		NewGetFlowChartHandler[adapters.WagtailDataModel](agr),
	}
}
