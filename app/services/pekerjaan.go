package services

import (
	"context"
	"latihan_uts_2/app/models"
	"latihan_uts_2/app/repository"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// IUserService mendefinisikan operasi logika bisnis untuk entitas User.
type IPekerjaanService interface {
    CreatePekerjaan(c *fiber.Ctx) error
    GetAllPekerjaan(c *fiber.Ctx) error
    GetPekerjaanByID(c *fiber.Ctx) error
    UpdatePekerjaan(c *fiber.Ctx) error
	SoftDeletePekerjaan(c *fiber.Ctx) error
}

// UserService implementasi IAlumniService.
type PekerjaanService struct {
    repo repository.IPekerjaanRepository // Ketergantungan pada Repository
}

// NewUserService membuat instance baru dari UserService.
func NewPekerjaanService(repo repository.IPekerjaanRepository) IPekerjaanService {
    return &PekerjaanService{repo: repo}
}

// CreateUser memvalidasi data dan meneruskannya ke repository.
func (s *PekerjaanService) CreatePekerjaan(c *fiber.Ctx) error {
	
	pekerjaan := new(models.Pekerjaan)
	if err := c.BodyParser(pekerjaan); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Permintaan tidak valid",
			"message": "Pastikan body permintaan Anda dalam format JSON yang benar.",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	now := time.Now()
	pekerjaan.CreatedAt = now
	pekerjaan.UpdatedAt = now

	if pekerjaan.NamaPerusahaan == "" {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error":   "Gagal membuat pengguna",
			"message": "nama tidak boleh kosong",
		})
	}
	if pekerjaan.BidangIndustri == "" || pekerjaan.PosisiJabatan <= "" {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error":   "Gagal membuat pengguna",
			"message": "email dan jurusan harus diisi dengan benar",
		})
	}

	createdPekerjaan, err := s.repo.CreatePekerjaan(ctx, pekerjaan)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error":   "Gagal membuat pengguna",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(createdPekerjaan)
}

// GetAllPekerjaan mengambil semua pengguna.
func (s *PekerjaanService) GetAllPekerjaan(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pekerjaan, err := s.repo.FindAllPekerjaan(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Gagal mengambil data",
			"message": err.Error(),
		})
	}

	return c.JSON(pekerjaan)
}

// GetPekerjaanByID mengambil pengguna
func (s *PekerjaanService) GetPekerjaanByID(c *fiber.Ctx) error {
	id := c.Params("id")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pekerjaan, err := s.repo.FindPekerjaanByID(ctx, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "Gagal mengambil pengguna",
			"message": err.Error(),
		})
	}
	if pekerjaan == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "Gagal mengambil pengguna",
			"message": "pengguna dengan ID tersebut tidak ditemukan",
		})
	}

	return c.JSON(pekerjaan)
}

// UpdatePekerjaan menangani logika pembaruan data.
func (s *PekerjaanService) UpdatePekerjaan(c *fiber.Ctx) error {
	id := c.Params("id")

	pekerjaan := new(models.Pekerjaan)
	pekerjaan.CreatedAt = time.Now()
	if err := c.BodyParser(pekerjaan); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Permintaan tidak valid",
			"message": "Pastikan body permintaan Anda dalam format JSON yang benar.",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pekerjaan.ID = primitive.NilObjectID
	pekerjaan.UpdatedAt = time.Now()

	if pekerjaan.NamaPerusahaan == "" {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error":   "Gagal memperbarui",
			"message": "nama perusahaan tidak boleh kosong saat diperbarui",
		})
	}

	updatedPekerjaan, err := s.repo.UpdatePekerjaan(ctx, id, pekerjaan)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   "Gagal memperbarui",
				"message": "pekerjaan dengan ID tersebut tidak ditemukan atau sudah dihapus",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Gagal memperbarui",
			"message": err.Error(),
		})
	}

	return c.JSON(updatedPekerjaan)
}

// SoftDeletePekerjaan
func (s *PekerjaanService) SoftDeletePekerjaan(c *fiber.Ctx) error {
	id := c.Params("id")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := s.repo.SoftDeletePekerjaan(ctx, id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   "Gagal menghapus",
				"message": "pekerjaan dengan ID tersebut tidak ditemukan atau sudah dihapus",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Gagal menghapus",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Pekerjaan berhasil dihapus (soft delete)",
	})
}