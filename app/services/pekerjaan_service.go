package services

import (
	"database/sql"
	"fmt"
	"latihan_uts_2/app/models"
	"latihan_uts_2/app/repository"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func GetAllPekerjaanService(c *fiber.Ctx, db *sql.DB) error {
	// Log siapa yang mengakses 
    username := c.Locals("username").(string) 
    log.Printf("User %s mengakses GET /api/pekerjaan", username)

	pekerjaan, err := repository.GetAllPekerjaan(db)
		if err == sql.ErrNoRows{
			return	c.Status(fiber.StatusOK).JSON(fiber.Map{
				"message":"Data Pekerjaan tidak ditemukan",
				"success": true,
				"isPekerjaan":	false,
			})
		}

	return c.Status(fiber.StatusOK).JSON(pekerjaan)
}

func GetPekerjaanByIDService(c *fiber.Ctx, db *sql.DB, id string) error {
	username := c.Locals("username").(string) 
    pekerjaanID, err := strconv.Atoi(id)
    if err != nil { 
        return c.Status(400).JSON(fiber.Map{ 
            "error": "ID tidak valid", 
        }) 
    }

	log.Printf("User %s mengakses GET /api/pekerjaan/%d", username, pekerjaanID) 
	
	pekerjaan, err := repository.GetPekerjaanByID(db, id)
	if err != nil{
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"pesan": "data tidak ditemukan",
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(pekerjaan)
}

func CreatePekerjaanService(c *fiber.Ctx, db *sql.DB) error{
	username := c.Locals("username").(string) 
    log.Printf("Admin %s menambah alumni baru", username) 

	var pekerjaan models.CreatePekerjaan

	if err := c.BodyParser(&pekerjaan); err != nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error" : err,
		})
	}
	CreatedPekerjaan, err := repository.CreatePekerjaan(db, pekerjaan)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"pesan": "Gagal menambahkan data",
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(CreatedPekerjaan)
}

func UpdatePekerjaanService(c *fiber.Ctx, db *sql.DB, id string) error {
	username := c.Locals("username").(string) 
    pekerjaanID, err := strconv.Atoi(id) 
    if err != nil { 
        return c.Status(400).JSON(fiber.Map{ 
            "error": "ID tidak valid", 
        }) 
    } 
 
    log.Printf("Admin %s mengupdate pekerjaan ID %d", username, pekerjaanID)

	var pekerjaan models.UpdatePekerjaan
    
	if err := c.BodyParser(&pekerjaan); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"pesan": "Format input tidak valid",
			"error": err.Error(),
		})
	}

	UpdatedPekerjaan, err := repository.UpdatePekerjaan(db, pekerjaan, id)
    
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"pesan": fmt.Sprintf("Data pekerjaan dengan ID %s tidak ditemukan untuk diperbarui.", id),
			})
		}
        
		fmt.Printf("Error saat memperbarui data pekerjaan ID %s: %v\n", id, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"pesan": "Error server: Gagal memperbarui data",
            // "detail_error": err.Error(), 
		})
	}

    // 4. Jika Sukses (data ditemukan dan diperbarui)
	return c.Status(fiber.StatusOK).JSON(UpdatedPekerjaan)
}

func DeletePekerjaanService(c *fiber.Ctx, db *sql.DB, id string) error {
	username := c.Locals("username").(string) 
    pekerjaanID, err := strconv.Atoi(id) 
    if err != nil { 
        return c.Status(400).JSON(fiber.Map{ 
            "error": "ID tidak valid", 
        }) 
    } 

	log.Printf("Admin %s menghapus alumni ID %d", username, pekerjaanID)

	var pekerjaan models.Pekerjaan
	if err := c.BodyParser(&pekerjaan); err != nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error" : err,
		})
	}

	DeletedPekerjaan, err := repository.DeletePekerjaan(db, pekerjaan , id)
	if err == sql.ErrNoRows {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"pesan": "data pekerjaan dengan ID " + id + " tidak ditemukan",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "pesan": "Data pekerjaan berhasil dihapus",
        "data_dihapus": DeletedPekerjaan,
    })
}