package routes

import (
	"database/sql"
	"latihan_uts_2/app/repository"
	"latihan_uts_2/app/services"
	"latihan_uts_2/middleware"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRoutes(app *fiber.App, mongoDB *mongo.Database, postgreSQL *sql.DB) {
	
	api := app.Group("/api")
	
    // Public routes (tidak perlu login)
	api.Group("") 
	api.Post("/login", func(c *fiber.Ctx) error {
		return Login(c, postgreSQL) 
	})

	// Protected routes (perlu login) 
	protected := api.Group("", middleware.AuthRequired()) 
    protected.Get("/profile", GetProfile) 

	// --- Grup Rute Alumni ---
    alumniRepo := repository.NewAlumniRepository(mongoDB)
    alumniService := services.NewAlumniService(alumniRepo)
	alumniRoutes := protected.Group("/alumni")
	SetupAlumniRoutes(alumniRoutes, alumniService)

	// --- Grup Rute Pekerjaan ---
	pekerjaanRepo := repository.NewPekerjaanRepository(mongoDB)
    pekerjaanService := services.NewPekerjaanService(pekerjaanRepo)
	pekerjaanRoutes := protected.Group("/pekerjaan")
	SetupPekerjaanRoutes(pekerjaanRoutes, pekerjaanService)

	// --- Grup Rute Trash ---
	trashed := protected.Group("/trash")
    trashed.Get("/", func(c *fiber.Ctx) error {
		return services.GetAllTrashService(c, postgreSQL)
	})

	trashed.Put("/:id", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		id := c.Params("id")
		return services.UpdateTrashService(c, postgreSQL, id)
	})

	trashed.Delete("/:id", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		id := c.Params("id")
		return services.DeleteTrashService(c, postgreSQL, id)
	})
}