package ports

import (
	"flowChart/handlers"
	"flowChart/transport"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type Encode struct {
	Success bool   `json:"success"`
	Err     string `json:"error"`
}

type HttpServer struct {
	App handlers.Application
}

func (h *HttpServer) EditFlowChartUnstructuredData(c *fiber.Ctx) error {
	ctx := c.Context()

	flowChartDto := &transport.FlowChartDto[transport.UnstructuredDataDto]{}

	if err := c.BodyParser(flowChartDto); err != nil {
		return c.Status(http.StatusBadRequest).JSON(Encode{Success: false, Err: err.Error()})
	}

	if err := h.App.Commands.EditFlowChart.Handler(ctx, flowChartDto); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(Encode{Success: false, Err: err.Error()})
	}

	return c.Status(http.StatusOK).JSON(Encode{Success: true, Err: ""})
}

func (h *HttpServer) GetFlowChartUnstructuredData(c *fiber.Ctx) error {
	ctx := c.Context()
	key := c.Params("key")

	flowChart, err := h.App.Queries.GetFlowChart.Handler(ctx, key)

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(Encode{Success: false, Err: err.Error()})
	}

	return c.Status(http.StatusOK).JSON(flowChart)
}
