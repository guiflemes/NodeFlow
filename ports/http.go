package ports

import (
	"flowChart/handlers"
	"flowChart/transport"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type Encode struct {
	Success bool  `json:"success"`
	Err     error `json:"error"`
}

type HttpServer struct {
	App handlers.Application
}

func (h *HttpServer) EditFlowChartSimpleData(c *fiber.Ctx) error {
	ctx := c.Context()

	flowChartDto := &transport.FlowChartDto[transport.DataDto]{}

	if err := c.BodyParser(flowChartDto); err != nil {
		return c.Status(http.StatusBadRequest).JSON(Encode{Success: false, Err: err})
	}

	if err := h.App.Commands.EditFlowChart.Handler(ctx, flowChartDto); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(Encode{Success: false, Err: err})
	}

	return c.Status(http.StatusOK).JSON(Encode{Success: true, Err: nil})
}
