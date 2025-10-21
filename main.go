package main

import (
	"context"
	"fmt"
	"latihan_uts_2/config"
	"latihan_uts_2/database"
	"latihan_uts_2/routes"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

func main() {
	
	// Menghubungkan .env
	config.Config()

	// database postgres
	postgresDB := database.ConnectDB()
	fmt.Println("Berhasil terhubung ke database PostgreSQL")

	// database mongodb
	mongoDB := database.ConnectMongoDB()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    defer func() {
        if err := mongoDB.Disconnect(ctx); err != nil {
            log.Fatalf("Gagal menutup koneksi MongoDB: %v", err)
        }
    }()

	// Inisialisasi fiber
	app := fiber.New(fiber.Config{
		ErrorHandler: func (c *fiber.Ctx, err error) error {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// routes
	routes.SetupRoutes(app, postgresDB, mongoDB)

	// Server
	log.Fatal(app.Listen(":3000"))
}