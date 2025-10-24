package main

import (
	"latihan_uts_2/config"
	"latihan_uts_2/database"
	"latihan_uts_2/routes"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	
	// Menghubungkan .env
	config.Config()

	// database postgresql
	postgreSQL := database.ConnectDB()

	// database mongodb
	mongoDB := database.ConnectMongoDB()

	// Inisialisasi fiber
	app := fiber.New(fiber.Config{
		ErrorHandler: func (c *fiber.Ctx, err error) error {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// routes
	routes.SetupRoutes(app, mongoDB, postgreSQL)

	// Server
	log.Fatal(app.Listen(":3000"))
}