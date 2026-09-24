package main

import (
	handler "herbie-charge-api/internal/charge/delivery/http"
	"herbie-charge-api/internal/charge/service"
	"herbie-charge-api/internal/config"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(
			".env not found, using system environment",
		)
	}

	cfg := config.LoadConfig()

	app := fiber.New()

	app.Use(
		cors.New(
			cors.Config{
				AllowOrigins:
					"http://localhost:4200," +
						"https://herbie-charge-angular.vercel.app",

				AllowHeaders:
					"Origin, Content-Type, Accept",
			},
		),
	)

	lineService :=
		service.NewLineService(
			cfg,
		)

	chargeService :=
		service.NewChargeService(
			lineService,
		)

	chargeHandler :=
		handler.NewChargeHandler(
			chargeService,
		)

	app.Get(
		"/health",
		func(c *fiber.Ctx) error {
			return c.JSON(
				fiber.Map{
					"status": "ok",
				},
			)
		},
	)

	app.Post(
		"/api/charge",
		chargeHandler.Charge,
	)

	log.Fatal(
		app.Listen(
			":" + cfg.Port,
		),
	)
}