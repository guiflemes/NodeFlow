package server

import (
	"flowChart/handlers"
	"flowChart/ports"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/sirupsen/logrus"
)

func RunHttpServer(addr string, application handlers.Application) {
	app := fiber.New()
	setMiddlewares(app)

	httpServer := ports.HttpServer{App: application}

	apiV1 := app.Group("api/v1")
	apiV1.Post("/flowchart", httpServer.EditFlowChartUnstructuredData)
	apiV1.Get("/flowchart/:key", httpServer.GetFlowChartUnstructuredData)

	logrus.Info("Starting HTTP server")
	app.Listen(addr)

}

func setMiddlewares(app *fiber.App) {
	addCorsMiddleware(app)
	addLoggingMiddleware(app)
}

func addCorsMiddleware(app *fiber.App) {
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))
}

func addLoggingMiddleware(app *fiber.App) {
	app.Use(logger.New(logger.Config{
		Next:         nil,
		Done:         nil,
		Format:       "[${time}] ${latency} | ${path} ${status} - ${method} \n",
		TimeFormat:   "02-Jan-2006",
		TimeZone:     "America/Sao_Paulo",
		TimeInterval: 500 * time.Millisecond,
		Output:       os.Stdout,
	}))
}
