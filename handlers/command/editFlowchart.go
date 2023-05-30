package command

import (
	"context"
	"flowChart/adapters"
	"flowChart/domain"
	"flowChart/transport"

	"fmt"
)

type dtoToDomain[R comparable, D comparable] func(flowChart *transport.FlowChartDto[R], dataParse func(request R) D) (*domain.FlowChart[D], error)
type dataParse[R comparable, D comparable] func(request R) D

type FlowCartRepo[T comparable] interface {
	StoreFlowChart(context.Context, *domain.FlowChart[T]) error
	UpdateFlowChart(context.Context, *domain.FlowChart[T]) error
	FlowChartExists(ctx context.Context, flowChart *domain.FlowChart[T]) (bool, error)
}

type EditHandlerFlowChart[R comparable, D comparable] struct {
	repo        FlowCartRepo[D]
	dtoToDomain dtoToDomain[R, D]
	parseData   dataParse[R, D]
}

func NewEditHandlerFlowChart[R comparable, D comparable](repo FlowCartRepo[D], parseData dataParse[R, D]) *EditHandlerFlowChart[R, D] {
	return &EditHandlerFlowChart[R, D]{
		repo:        repo,
		dtoToDomain: transport.ToDomain[R, D],
		parseData:   parseData,
	}
}

func (h *EditHandlerFlowChart[R, D]) Handler(ctx context.Context, dto *transport.FlowChartDto[R]) error {
	domain, err := h.dtoToDomain(dto, h.parseData)

	if err != nil {
		return fmt.Errorf("error parsing dto to domain %w", err)
	}

	exists, err := h.repo.FlowChartExists(ctx, domain)

	if err != nil {
		return err
	}

	if exists {
		return h.repo.UpdateFlowChart(ctx, domain)
	}

	return h.repo.StoreFlowChart(ctx, domain)

}

type HandlerFlowChartSimpleData struct {
	*EditHandlerFlowChart[transport.DataDto, domain.Data]
}

func NewHandlerFlowChartSimpleData(agr *adapters.FlowChartDataAggregate) HandlerFlowChartSimpleData {
	return HandlerFlowChartSimpleData{
		NewEditHandlerFlowChart[transport.DataDto, domain.Data](agr, func(request transport.DataDto) domain.Data { return domain.Data{Label: request.Label} }),
	}
}

type HandlerFlowChartUnstructuredData struct {
	*EditHandlerFlowChart[transport.UnstructuredDataDto, domain.UnstructuredDataDomain]
}

func NewHandlerFlowChartUnstructuredData(agr *adapters.WriteFlowChartUnstructuredDataAgg) HandlerFlowChartUnstructuredData {
	return HandlerFlowChartUnstructuredData{
		NewEditHandlerFlowChart[transport.UnstructuredDataDto, domain.UnstructuredDataDomain](agr,
			func(request transport.UnstructuredDataDto) domain.UnstructuredDataDomain {
				return request
			}),
	}
}
