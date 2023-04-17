package server

import (
	"flowChart/handlers"
	"flowChart/ports"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/sirupsen/logrus"
)

func RunHttpServer(addr string, application handlers.Application) {
	app := fiber.New()
	setMiddlewares(app)

	httpServer := ports.HttpServer{App: application}

	apiV1 := app.Group("api/v1")
	apiV1.Post("/flowchart", httpServer.EditFlowChartSimpleData)

	logrus.Info("Starting HTTP server")

	if err := app.Listen(addr); err != nil {
		logrus.WithError(err).Panic("Unable to start HTTP server")
	}

}

func setMiddlewares(app *fiber.App) {}

func addCorsMiddleware(app *fiber.App) {
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))
}
