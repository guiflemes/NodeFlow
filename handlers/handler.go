package handlers

import (
	"context"
	"flowChart/domain"
	"flowChart/ports"
	"fmt"
)

type dtoToDomain[T comparable] func(flowChart *ports.FlowChartDto[T]) (*domain.FlowChart[T], error)

type FlowCartRepo[T comparable] interface {
	StoreFlowChart(context.Context, *domain.FlowChart[T]) error
	UpdateFlowChart(context.Context, *domain.FlowChart[T]) error
	FlowChartExists(ctx context.Context, flowChartKey string) (bool, error)
}

type HandlerFlowChart[T comparable] struct {
	repo        FlowCartRepo[T]
	dtoToDomain dtoToDomain[T]
}

func (h *HandlerFlowChart[T]) Handler(ctx context.Context, dto *ports.FlowChartDto[T]) error {
	domain, err := h.dtoToDomain(dto)

	if err != nil {
		return fmt.Errorf("error parsing dto to domain %w", err)
	}

	exists, err := h.repo.FlowChartExists(ctx, dto.Key)

	if err != nil {
		return err
	}

	if exists {
		return h.repo.UpdateFlowChart(ctx, domain)
	}

	return h.repo.StoreFlowChart(ctx, domain)

}
