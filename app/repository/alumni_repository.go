package repository

import (
	"database/sql"
	"fmt"
	"latihan_uts_2/app/models"
)

func GetAllAlumni(search, sortBy, order string, limit, offset int, db *sql.DB) ([]models.Alumni,error) {
	query := fmt.Sprintf(` 
		SELECT id, nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, created_at, updated_at
		FROM alumni
		WHERE nim ILIKE $1 OR nama ILIKE $1
		ORDER BY %s %s
		LIMIT $2 OFFSET $3
    `, sortBy, order) 

	rows, err := db.Query(query,  "%"+search+"%", limit, offset)

	if err != nil{
		return nil, err
	}

	defer rows.Close()

	var alumniList  []models.Alumni

	for rows.Next(){
		var a models.Alumni
		err := rows.Scan(
			&a.ID, 
			&a.NIM, 
			&a.Nama, 
			&a.Jurusan, 
			&a.Angkatan, 
			&a.TahunLulus, 
			&a.Email, 
			&a.NoTelepon, 
			&a.Alamat, 
			&a.CreatedAt, 
			&a.UpdatedAt,
		)
		
	if err != nil{
		if err == sql.ErrNoRows {
			return []models.Alumni{}, sql.ErrNoRows 
		}

		fmt.Println("Error saat memindai data pekerjaan:", err) 
		return []models.Alumni{}, err
	}

		alumniList = append(alumniList, a)
	}

	return alumniList, err
}

// CountUsersRepo -> hitung total data untuk pagination  FOR PAGINATION
func CountAlumniRepo(search string, db *sql.DB) (int, error) { 
    var total int 
    countQuery := `SELECT COUNT(*) FROM alumni WHERE nim ILIKE $1 OR nama ILIKE $1` 
    err := db.QueryRow(countQuery, "%"+search+"%").Scan(&total) 
    if err != nil && err != sql.ErrNoRows { 
        return 0, err 
    } 
    return total, nil 
} 

func GetAlumniByID(db *sql.DB, id string) (models.Alumni, error){
	var alumni models.Alumni

	query := `
		SELECT id, nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, created_at, updated_at
		FROM alumni
		WHERE id = $1
	`
	row := db.QueryRow(query, id)
	err := row.Scan(
		&alumni.ID,
		&alumni.NIM,
		&alumni.Nama,
		&alumni.Jurusan,
		&alumni.Angkatan,
		&alumni.TahunLulus,
		&alumni.Email,
		&alumni.NoTelepon,
		&alumni.Alamat,
		&alumni.CreatedAt,
		&alumni.UpdatedAt,
	)
	if err != nil{
		return alumni, sql.ErrNoRows
	}
	
	return alumni, nil
}

func CreateAlumni(db *sql.DB, alumni models.CreateAlumni) (models.CreateAlumni, error){
	query := `
		INSERT INTO alumni (nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat
	`
	var created models.CreateAlumni

	err := db.QueryRow(query,
		&alumni.NIM,
		&alumni.Nama,
		&alumni.Jurusan,
		&alumni.Angkatan,
		&alumni.TahunLulus,
		&alumni.Email,
		&alumni.NoTelepon,
		&alumni.Alamat,
	).Scan(
		&created.NIM,
		&created.Nama,
		&created.Jurusan,
		&created.Angkatan,
		&created.TahunLulus,
		&created.Email,
		&created.NoTelepon,
		&created.Alamat,
	)
	if err != nil{
		return created, err
	}

	return created, nil
}

func UpdateAlumni(db *sql.DB, alumni models.UpdateAlumni, id string) (models.UpdateAlumni, error){
	query := `
		UPDATE alumni
		SET nama = $1, email = $2, no_telepon = $3, alamat = $4
		WHERE id = $5
		RETURNING nama, email, no_telepon, alamat
	`
	var updated models.UpdateAlumni

	err := db.QueryRow(query,
		&alumni.Nama,
		&alumni.Email,
		&alumni.NoTelepon,
		&alumni.Alamat,
		id,
	).Scan(
		&updated.Nama,
		&updated.Email,
		&updated.NoTelepon,
		&updated.Alamat,
	)
	if err != nil{
		return updated, err
	}

	return updated, nil
}

func DeleteAlumni(db *sql.DB, id string) (models.Alumni, error) {
	var alumni models.Alumni

	query := `
		DELETE FROM alumni
		WHERE id = $1
		RETURNING id, nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat
	`
	    err := db.QueryRow(query, id).Scan(
        &alumni.ID,
        &alumni.NIM,
        &alumni.Nama,
        &alumni.Jurusan,
        &alumni.Angkatan,
        &alumni.TahunLulus,
        &alumni.Email,
        &alumni.NoTelepon,
        &alumni.Alamat,
    )

    return alumni, err
}

func SoftDeleteAlumni(db *sql.DB, id string) (models.Alumni, error) {
	var alumni models.Alumni

	query := `
		UPDATE alumni
		SET deleted_at = NOW()
		WHERE id = $1
		RETURNING id, nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, deleted_at, created_at, updated_at
	`

	err := db.QueryRow(query, id).Scan(
		&alumni.ID,
		&alumni.NIM,
		&alumni.Nama,
		&alumni.Jurusan,
		&alumni.Angkatan,
		&alumni.TahunLulus,
		&alumni.Email,
		&alumni.NoTelepon,
		&alumni.Alamat,
        
		&alumni.DeletedAt,
		&alumni.CreatedAt,
		&alumni.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.Alumni{}, sql.ErrNoRows
		}
		return models.Alumni{}, err
	}

	return alumni, nil
}