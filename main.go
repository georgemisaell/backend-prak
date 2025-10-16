package main

import (
	"fmt"
	"latihan_uts_2/config"
	"latihan_uts_2/database"
	"latihan_uts_2/routes"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	
	// Menghubungkan .env
	config.Config()

	// database
	db := database.ConnectDB()
	fmt.Println("Berhasil terhubung ke database PostgreSQL")

	// Inisialisasi fiber
	app := fiber.New(fiber.Config{
		ErrorHandler: func (c *fiber.Ctx, err error) error {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// routes
	routes.SetupRoutes(app, db)

	// Server
	log.Fatal(app.Listen(":3000"))
}