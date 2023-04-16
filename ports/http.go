package ports

import (
	"github.com/gofiber/fiber/v2"
)

func CreateFlowChartData(c *fiber.Ctx) error {
	c.Context()

	flowChartDto := &FlowChartDto[string]{}

	if err := c.BodyParser(flowChartDto); err != nil {
		return err
	}

	return nil
}
