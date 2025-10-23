package routes

import (
	"context"
	"latihan_uts_2/app/models"
	"latihan_uts_2/app/services"
	"time"

	"github.com/gofiber/fiber/v2"
)

func SetupAlumniRoutes(alumniRoutes fiber.Router, alumniService services.IAlumniService) {
	// POST /api/v1/users -> Membuat pengguna baru
	alumniRoutes.Post("/", func(c *fiber.Ctx) error {
		user := new(models.Alumni)
		if err := c.BodyParser(user); err != nil {
			// Permintaan buruk (misalnya, JSON tidak valid)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Permintaan tidak valid",
				"message": "Pastikan body permintaan Anda dalam format JSON yang benar.",
			})
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		createdAlumni, err := alumniService.CreateAlumni(ctx, user)
		if err != nil {
			// Menangani error dari layer Service/Repository (misalnya, validasi gagal, DB error)
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"error":   "Gagal membuat pengguna",
				"message": err.Error(),
			})
		}

		return c.Status(fiber.StatusCreated).JSON(createdAlumni)
	})

	// GET /api/v1/users -> Mendapatkan semua pengguna
	alumniRoutes.Get("/", func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		users, err := alumniService.GetAllAlumni(ctx)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Gagal mengambil data",
				"message": err.Error(),
			})
		}
		return c.JSON(users)
	})

	    // GET /api/v1/users/:id -> Mendapatkan pengguna berdasarkan ID
    alumniRoutes.Get("/:id", func(c *fiber.Ctx) error {
        id := c.Params("id")

        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()

        user, err := alumniService.GetAlumniByID(ctx, id)
        if err != nil {
            // Error dari service menunjukkan pengguna tidak ditemukan atau ID tidak valid
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
                "error":   "Gagal mengambil pengguna",
                "message": err.Error(),
            })
        }

        return c.JSON(user)
    })
}