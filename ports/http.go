package ports

import (
	"flowChart/handlers"
	"flowChart/transport"

	"github.com/gofiber/fiber/v2"
)

type HttpServer struct {
	App handlers.Application
}

func (h *HttpServer) EditFlowChartSimpleData(c *fiber.Ctx) error {
	ctx := c.Context()

	flowChartDto := &transport.FlowChartDto[transport.DataDto]{}

	if err := c.BodyParser(flowChartDto); err != nil {
		return err
	}

	if err := h.App.Commands.EditFlowChart.Handler(ctx, flowChartDto); err != nil {
		return err
	}

	return nil
}
