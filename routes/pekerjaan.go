package routes

import (
	"latihan_uts_2/app/services"
	"latihan_uts_2/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupPekerjaanRoutes(pekerjaanRoutes fiber.Router, pekerjaanService services.IPekerjaanService) {

	// POST /api1/v1/pekerjaan -> Membuat pekerjaan baru
	pekerjaanRoutes.Post("/", func(c *fiber.Ctx) error {
		// Semua logika dipindahkan ke service
		return pekerjaanService.CreatePekerjaan(c)
	})

	// GET /api1/v1/pekerjaan -> Mendapatkan semua pekerjaan
	pekerjaanRoutes.Get("/", func(c *fiber.Ctx) error {
		// Semua logika dipindahkan ke service
		return pekerjaanService.GetAllPekerjaan(c)
	})

	// GET /api1/v1/pekerjaan/:id -> Mendapatkan pekerjaan berdasarkan ID
	pekerjaanRoutes.Get("/:id", func(c *fiber.Ctx) error {
		// Semua logika dipindahkan ke service
		return pekerjaanService.GetPekerjaanByID(c)
	})

	// PUT /api1/v1/pekerjaan/:id -> Update pekerjaan
	pekerjaanRoutes.Put("/:id", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return pekerjaanService.UpdatePekerjaan(c)
	})

	// DELETE /api1/v1/pekerjaan/:id -> Soft delete pekerjaan
	pekerjaanRoutes.Delete("/:id", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return pekerjaanService.SoftDeletePekerjaan(c)
	})
}