package routes

import (
	"latihan_uts_2/app/services"
	"latihan_uts_2/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetFileUploadRoutes(filesRoutes fiber.Router, filesService services.FileService) {
	
	// POST /api1/v1/file-upload -> Membuat file-upload baru
	filesRoutes.Post("/upload", middleware.AdminOnly(),filesService.UploadFile) 
	
	// GET /api1/v1/file-upload -> Mendapatkan semua file-upload
	filesRoutes.Get("/", middleware.AdminOnly(),filesService.GetAllFiles) 
	
	// GET /api1/v1/file-upload/:id -> Mendapatkan file-upload berdasarkan ID
	filesRoutes.Get("/:id", middleware.AdminOnly(),filesService.GetFileByID) 
	
	// DELETE /api1/v1/file-upload/:id -> Soft delete file-upload
	filesRoutes.Delete("/:id", middleware.AdminOnly(),filesService.DeleteFile)

}