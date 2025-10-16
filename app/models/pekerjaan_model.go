package models

import (
	"time"
)

type Pekerjaan struct {
	ID                int    `json:"id"`
	AlumniID          int    `json:"alumni_id"`
	NamaPerusahaan    string `json:"nama_perusahaan"`
	PosisiJabatan     string `json:"posisi_jabatan"`
	BidangIndustri    string `json:"bidang_industri"`
	LokasiKerja       string `json:"lokasi_kerja"`
	GajiRange         string `json:"gaji_range"`
	TanggalMulaiKerja time.Time `json:"tanggal_mulai_kerja"`
	TanggalSelesaiKerja *time.Time `json:"tanggal_selesai_kerja"`
	StatusPekerjaan string `json:"status_pekerjaan"`
	DeskripsiPekerjaan *string `json:"deskripsi_pekerjaan"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreatePekerjaan struct{
	AlumniID          int    `json:"alumni_id"`
	NamaPerusahaan    string `json:"nama_perusahaan"`
	PosisiJabatan     string `json:"posisi_jabatan"`
	BidangIndustri    string `json:"bidang_industri"`
	LokasiKerja       string `json:"lokasi_kerja"`
	GajiRange         string `json:"gaji_range"`
	TanggalMulaiKerja time.Time `json:"tanggal_mulai_kerja"`
	TanggalSelesaiKerja time.Time `json:"tanggal_selesai_kerja"`
	StatusPekerjaan string `json:"status_pekerjaan"`
}

type UpdatePekerjaan struct{
	NamaPerusahaan    string `json:"nama_perusahaan"`
	PosisiJabatan     string `json:"posisi_jabatan"`
	BidangIndustri    string `json:"bidang_industri"`
	LokasiKerja       string `json:"lokasi_kerja"`
	GajiRange         string `json:"gaji_range"`
	StatusPekerjaan string `json:"status_pekerjaan"`
	DeskripsiPekerjaan string `json:"deskripsi_pekerjaan"`
	TanggalMulaiKerja time.Time `json:"tanggal_mulai_kerja"`
	TanggalSelesaiKerja time.Time `json:"tanggal_selesai_kerja"`
}

type PekerjaanResponse struct { 
    Data []Pekerjaan   `json:"data"` 
    Meta MetaInfo `json:"meta"` 
} 