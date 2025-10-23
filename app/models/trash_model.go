package models

import "time"

type Trash struct {
	ID                int    `bson:"id" json:"id"`
	AlumniID          int    `bson:"alumni_id" json:"alumni_id"`
	NamaPerusahaan    string `bson:"nama_perusahaan" json:"nama_perusahaan"`
	PosisiJabatan     string `bson:"posisi_jabatan" json:"posisi_jabatan"`
	BidangIndustri    string `bson:"bidang_industri" json:"bidang_industri"`
	LokasiKerja       string `bson:"lokasi_kerja" json:"lokasi_kerja"`
	GajiRange         string `bson:"gaji_range" json:"gaji_range"`
	TanggalMulaiKerja time.Time `bson:"tanggal_mulai_kerja" json:"tanggal_mulai_kerja"`
	TanggalSelesaiKerja *time.Time `bson:"tanggal_selesai_kerja" json:"tanggal_selesai_kerja"`
	StatusPekerjaan string `bson:"status_pekerjaan" json:"status_pekerjaan"`
	DeskripsiPekerjaan *string `bson:"deskripsi_pekerjaan" json:"deskripsi_pekerjaan"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
	DeletedAt  *time.Time `bson:"deleted_at" json:"deleted_at"`
}

type UpdateTrash struct{
	DeletedAt  *time.Time `bson:"deleted_at" json:"deleted_at"`
}

type TrashResponse struct { 
    Data []Trash   `bson:"data" json:"data"` 
    Meta MetaInfo `bson:"meta" json:"meta"` 
} 