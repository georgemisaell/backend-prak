package repository

import (
	"database/sql"
	"fmt"
	"latihan_uts_2/app/models"
)

func GetAllTrash(search, sortBy, order string, limit, offset int, db *sql.DB) ([]models.Trash, error) {
	var trash []models.Trash

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
			updated_at,
			deleted_at
		FROM pekerjaan_alumni
		WHERE deleted_at IS NOT NULL
		  AND (nama_perusahaan ILIKE $1 OR bidang_industri ILIKE $1)
		ORDER BY %s %s
		LIMIT $2 OFFSET $3
	`, sortBy, order)

	rows, err := db.Query(query, "%"+search+"%", limit, offset)
	if err != nil {
		fmt.Println("Gagal mengambil seluruh data sampah:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var row models.Trash
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
			&row.DeletedAt,
		)
		if err != nil {
			fmt.Println("Error saat memindai data sampah:", err)
			return nil, err
		}

		trash = append(trash, row)
	}

	return trash, nil
}

func UpdateTrash(db *sql.DB, trash models.UpdateTrash, id string) (models.UpdateTrash, error) {
	query := `
		UPDATE pekerjaan_alumni
		SET deleted_at = NULL
		WHERE id = $1
		RETURNING deleted_at
	`

	var trashUpdated models.UpdateTrash
	err := db.QueryRow(query, id).Scan(
		&trashUpdated.DeletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.UpdateTrash{}, sql.ErrNoRows
		}

		fmt.Println("Error saat memindai data hasil update:", err)
		return models.UpdateTrash{}, err
	}

	return trashUpdated, nil
}

func DeleteTrash(db *sql.DB, trash models.Trash, id string) (models.Trash, error){
	query := `
		DELETE FROM pekerjaan_alumni
		WHERE id = $1
		RETURNING id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at
	`

	err := db.QueryRow(query, id).Scan(
		&trash.ID,
		&trash.AlumniID,
		&trash.NamaPerusahaan,
		&trash.PosisiJabatan,
		&trash.BidangIndustri,
		&trash.LokasiKerja,
		&trash.GajiRange,
		&trash.TanggalMulaiKerja,
		&trash.TanggalSelesaiKerja,
		&trash.StatusPekerjaan,
		&trash.DeskripsiPekerjaan,
		&trash.CreatedAt,
		&trash.UpdatedAt,
	)
	if err != nil{
		return trash, err
	}

	return trash, err
}