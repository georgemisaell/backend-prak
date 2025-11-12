package services

import (
	"context"
	"errors"
	"latihan_uts/app/models"
	"latihan_uts/app/repository"
	"time"
)

// IUserService mendefinisikan operasi logika bisnis untuk entitas User.
type IAlumniService interface {
    CreateAlumni(ctx context.Context, user *models.Alumni) (*models.Alumni, error)
    GetAlumniByID(ctx context.Context, id string) (*models.Alumni, error)
    GetAllAlumni(ctx context.Context) ([]models.Alumni, error)
}

// UserService implementasi IAlumniService.
type AlumniService struct {
    repo repository.IAlumniRepository // Ketergantungan pada Repository
}

// NewUserService membuat instance baru dari UserService.
func NewAlumniService(repo repository.IAlumniRepository) IAlumniService {
    return &AlumniService{repo: repo}
}

// CreateUser memvalidasi data dan meneruskannya ke repository.
// @Summary      Membuat data alumni baru
// @Description  Menambahkan alumni ke database
// @Tags         alumni
// @Accept       json
// @Produce      json
// @Param        alumni  body      models.Alumni  true  "Data Alumni"
// @Router       /alumni [post]
func (s *AlumniService) CreateAlumni(ctx context.Context, alumni *models.Alumni) (*models.Alumni, error) {

    now := time.Now()
    alumni.CreatedAt = now
    alumni.UpdatedAt = now

    if alumni.Nama == "" {
        return nil, errors.New("nama tidak boleh kosong")
    }	
    if alumni.Email == "" || alumni.Jurusan <= "" {
        return nil, errors.New("email dan jurusan harus diisi dengan benar")
    }

    return s.repo.CreateAlumni(ctx, alumni)
}

// GetUserByID mengambil pengguna dan menangani kasus jika tidak ditemukan.
// GetAlumniByID godoc
// @Summary Dapatkan alumni berdasarkan id
// @Tags Users - Alumni
// @Router /alumni/{id} [get]
func (s *AlumniService) GetAlumniByID(ctx context.Context, id string) (*models.Alumni, error) {
    alumni, err := s.repo.FindAlumniByID(ctx, id)
    if err != nil {
        return nil, err
    }
    if alumni == nil {
        return nil, errors.New("pengguna dengan ID tersebut tidak ditemukan")
    }
    return alumni, nil
}

// GetAllUsers mengambil semua pengguna.
// GetAllAlumni godoc
// @Summary Dapatkan semua alumni
// @Tags Users - Alumni
// @Router /alumni [get]
func (s *AlumniService) GetAllAlumni(ctx context.Context) ([]models.Alumni, error) {
    alumni, err := s.repo.FindAllAlumni(ctx)
    if err != nil {
        return nil, err
    }
    return alumni, nil
}