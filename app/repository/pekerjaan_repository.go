package repository

import (
	"database/sql"
	"fmt"
	"latihan_uts_2/app/models"
)

func GetAllPekerjaan(search, sortBy, order string, limit, offset int, db *sql.DB)([]models.Pekerjaan, error){
	var pekerjaan []models.Pekerjaan

	query := fmt.Sprintf(`
		SELECT 
			id,
			alumni_id,
			nama_perusahaan,
			posisi_jabatan,
			bidang_industri,
			lokasi_kerja,
			gaji_range,
			tanggal_mulai_kerja,
			tanggal_selesai_kerja,
			status_pekerjaan,
			deskripsi_pekerjaan,
			created_at,
			updated_at
		FROM pekerjaan_alumni
		WHERE nama_perusahaan ILIKE $1 OR bidang_industri ILIKE $1
		ORDER BY %s %s
		LIMIT $2 OFFSET $3
	`, sortBy, order)

	rows, err := db.Query(query,  "%"+search+"%", limit, offset)
	if err != nil {
		fmt.Println("Gagal mengambil seluruh pekerjaan alumni")
	}

	for rows.Next(){
		var row models.Pekerjaan
		err := rows.Scan(
			&row.ID,
			&row.AlumniID,
			&row.NamaPerusahaan,
			&row.PosisiJabatan,
			&row.BidangIndustri,
			&row.LokasiKerja,
			&row.GajiRange,
			&row.TanggalMulaiKerja,
			&row.TanggalSelesaiKerja,
			&row.StatusPekerjaan,
			&row.DeskripsiPekerjaan,
			&row.CreatedAt,
			&row.UpdatedAt,
		)

	if err != nil{
		if err == sql.ErrNoRows {
			return []models.Pekerjaan{}, sql.ErrNoRows 
		}

		fmt.Println("Error saat memindai data pekerjaan:", err) 
		return []models.Pekerjaan{}, err
	}

		pekerjaan = append(pekerjaan, row)
	}
	
	return pekerjaan, err

}

// CountUsersRepo -> hitung total data untuk pagination  FOR PAGINATION
func CountPekerjaanRepo(search string, db *sql.DB) (int, error) { 
    var total int 
    countQuery := `SELECT COUNT(*) FROM alumni WHERE nama_perusahaan ILIKE $1 OR bidang_industri ILIKE $1` 
    err := db.QueryRow(countQuery, "%"+search+"%").Scan(&total) 
    if err != nil && err != sql.ErrNoRows { 
        return 0, err 
    } 
    return total, nil 
} 

func GetPekerjaanByID(db *sql.DB, id string)(models.Pekerjaan, error){
	query := `
		SELECT 
			id,
			alumni_id,
			nama_perusahaan,
			posisi_jabatan,
			bidang_industri,
			lokasi_kerja,
			gaji_range,
			tanggal_mulai_kerja,
			tanggal_selesai_kerja,
			status_pekerjaan,
			deskripsi_pekerjaan,
			created_at,
			updated_at
		FROM pekerjaan_alumni
		WHERE id = $1
	`

	rows := db.QueryRow(query, id)

	var pekerjaan models.Pekerjaan
	err := rows.Scan(
		&pekerjaan.ID,
		&pekerjaan.AlumniID,
		&pekerjaan.NamaPerusahaan,
		&pekerjaan.PosisiJabatan,
		&pekerjaan.BidangIndustri,
		&pekerjaan.LokasiKerja,
		&pekerjaan.GajiRange,
		&pekerjaan.TanggalMulaiKerja,
		&pekerjaan.TanggalSelesaiKerja,
		&pekerjaan.StatusPekerjaan,
		&pekerjaan.DeskripsiPekerjaan,
		&pekerjaan.CreatedAt,
		&pekerjaan.UpdatedAt,
	)
	
	if err != nil{
		if err == sql.ErrNoRows {
			return models.Pekerjaan{}, sql.ErrNoRows 
		}

		fmt.Println("Error saat memindai data pekerjaan:", err) 
		return models.Pekerjaan{}, err
	}

	return pekerjaan, err
}

func CreatePekerjaan(db *sql.DB, pekerjaan models.CreatePekerjaan)(models.CreatePekerjaan, error){
	query := `
		INSERT INTO pekerjaan_alumni (alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan
	`

	var Created models.CreatePekerjaan
	err := db.QueryRow(query, 
		pekerjaan.AlumniID,
		pekerjaan.NamaPerusahaan,
		pekerjaan.PosisiJabatan,
		pekerjaan.BidangIndustri,
		pekerjaan.LokasiKerja,
		pekerjaan.GajiRange,
		pekerjaan.TanggalMulaiKerja,
		pekerjaan.TanggalSelesaiKerja,
		pekerjaan.StatusPekerjaan,
	).Scan(
		&Created.AlumniID,
		&Created.NamaPerusahaan,
		&Created.PosisiJabatan,
		&Created.BidangIndustri,
		&Created.LokasiKerja,
		&Created.GajiRange,
		&Created.TanggalMulaiKerja,
		&Created.TanggalSelesaiKerja,
		&Created.StatusPekerjaan,
	)

	if err != nil {
		fmt.Println("Gagal menambahkan pekerjaan ", err)
	}

	return Created, err
}

func UpdatePekerjaan(db *sql.DB, pekerjaan models.UpdatePekerjaan, id string)(models.UpdatePekerjaan, error){

	query := `
		UPDATE pekerjaan_alumni
		SET nama_perusahaan=$1, posisi_jabatan= $2, bidang_industri=$3, lokasi_kerja=$4, gaji_range=$5, status_pekerjaan=$6, deskripsi_pekerjaan=$7 , tanggal_mulai_kerja=$8 , tanggal_selesai_kerja=$9
		WHERE id = $10
		RETURNING nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range, status_pekerjaan, deskripsi_pekerjaan, tanggal_mulai_kerja, tanggal_selesai_kerja
	`
	var pekerjaanUpdated models.UpdatePekerjaan
	err := db.QueryRow(query,
		&pekerjaan.NamaPerusahaan,
		&pekerjaan.PosisiJabatan,
		&pekerjaan.BidangIndustri,
		&pekerjaan.LokasiKerja,
		&pekerjaan.GajiRange,
		&pekerjaan.StatusPekerjaan,
		&pekerjaan.DeskripsiPekerjaan,
		&pekerjaan.TanggalMulaiKerja,
		&pekerjaan.TanggalSelesaiKerja,
		id,
	).Scan(
		&pekerjaanUpdated.NamaPerusahaan,
		&pekerjaanUpdated.PosisiJabatan,
		&pekerjaanUpdated.BidangIndustri,
		&pekerjaanUpdated.LokasiKerja,
		&pekerjaanUpdated.GajiRange,
		&pekerjaanUpdated.StatusPekerjaan,
		&pekerjaan.DeskripsiPekerjaan,
		&pekerjaanUpdated.TanggalMulaiKerja,
		&pekerjaanUpdated.TanggalSelesaiKerja,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.UpdatePekerjaan{}, sql.ErrNoRows 
		}
        
		fmt.Println("Error saat memindai data hasil update:", err)
		return models.UpdatePekerjaan{}, err
	}

	return pekerjaanUpdated, nil
}

func DeletePekerjaan(db *sql.DB, pekerjaan models.Pekerjaan, id string) (models.Pekerjaan, error){
	query := `
		DELETE FROM pekerjaan_alumni
		WHERE id = $1
		RETURNING id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at
	`

	err := db.QueryRow(query, id).Scan(
		&pekerjaan.ID,
		&pekerjaan.AlumniID,
		&pekerjaan.NamaPerusahaan,
		&pekerjaan.PosisiJabatan,
		&pekerjaan.BidangIndustri,
		&pekerjaan.LokasiKerja,
		&pekerjaan.GajiRange,
		&pekerjaan.TanggalMulaiKerja,
		&pekerjaan.TanggalSelesaiKerja,
		&pekerjaan.StatusPekerjaan,
		&pekerjaan.DeskripsiPekerjaan,
		&pekerjaan.CreatedAt,
		&pekerjaan.UpdatedAt,
	)
	if err != nil{
		return pekerjaan, err
	}

	return pekerjaan, err
}

func SoftDeletePekerjaan(db *sql.DB, pekerjaan models.Pekerjaan, id string) (models.Pekerjaan, error) {

	query := `
		UPDATE pekerjaan_alumni
		SET deleted_at = NOW()
		WHERE id = $1
		RETURNING id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at, deleted_at
	`

	err := db.QueryRow(query, id).Scan(
		&pekerjaan.ID,
		&pekerjaan.AlumniID,
		&pekerjaan.NamaPerusahaan,
		&pekerjaan.PosisiJabatan,
		&pekerjaan.BidangIndustri,
		&pekerjaan.LokasiKerja,
		&pekerjaan.GajiRange,
		&pekerjaan.TanggalMulaiKerja,
		&pekerjaan.TanggalSelesaiKerja,
		&pekerjaan.StatusPekerjaan,
		&pekerjaan.DeskripsiPekerjaan,
		&pekerjaan.CreatedAt,
		&pekerjaan.UpdatedAt,        
		&pekerjaan.DeletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.Pekerjaan{}, sql.ErrNoRows
		}
		return models.Pekerjaan{}, err
	}

	return pekerjaan, nil
}