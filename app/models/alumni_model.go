package models

import "time"

type Alumni struct {
	ID         int     `json:"id"`
	NIM        string  `json:"nim"`
	Nama       string  `json:"nama"`
	Jurusan    string  `json:"jurusan"`
	Angkatan   int     `json:"angkatan"`
	TahunLulus int     `json:"tahun_lulus"`
	Email      string  `json:"email"`
	NoTelepon  *string `json:"no_telepon"`
	Alamat     string  `json:"alamat"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type CreateAlumni struct {
	NIM        string  `json:"nim"`
	Nama       string  `json:"nama"`
	Jurusan    string  `json:"jurusan"`
	Angkatan   int     `json:"angkatan"`
	TahunLulus int     `json:"tahun_lulus"`
	Email      string  `json:"email"`
	NoTelepon  *string `json:"no_telepon"`
	Alamat     string  `json:"alamat"`
}

type UpdateAlumni struct {
	Nama       string  `json:"nama"`
	Email      string  `json:"email"`
	NoTelepon  *string `json:"no_telepon"`
	Alamat     string  `json:"alamat"`
}

type AlumniResponse struct { 
    Data []Alumni   `json:"data"` 
    Meta MetaInfo `json:"meta"` 
} 