package main

import (
	"latihan_uts/config"
	"latihan_uts/database"
	_ "latihan_uts/docs"
	"latihan_uts/routes"
	"log"

	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

// @title Manajemen Alumni
// @version 1.0
// @description API untuk mengelola data dengan MongoDB menggunakan Clean Architecture
// @host 127.0.0.1:3000
// @BasePath /api
// @schemes http
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

	// Swagger UI
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// routes
	routes.SetupRoutes(app, mongoDB, postgreSQL)

	// Server
	log.Fatal(app.Listen(":3000"))
}