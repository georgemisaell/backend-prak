package services

import (
	"database/sql"
	"fmt"
	"latihan_uts_2/app/models"
	"latihan_uts_2/app/repository"
	"log"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func GetAllTrashService(c *fiber.Ctx, db *sql.DB) error {
	// Log siapa yang mengakses 
    username := c.Locals("username").(string) 
    log.Printf("User %s mengakses GET /api/trash", username)

	// Pagination
	page, _ := strconv.Atoi(c.Query("page", "1")) 
    limit, _ := strconv.Atoi(c.Query("limit", "10")) 
    sortBy := c.Query("sortBy", "id") 
    order := c.Query("order", "asc") 
    search := c.Query("search", "") 
	offset := (page - 1) * limit

	// Validasi input 
    sortByWhitelist := map[string]bool{"id": true, "nama_perusahaan": true, "bidang_industri": true, "created_at": true} 
    if !sortByWhitelist[sortBy] { 
        sortBy = "id" 
    } 
    if strings.ToLower(order) != "desc" { 
        order = "asc" 
    }

	trash, err := repository.GetAllTrash(search, sortBy, order, limit, offset, db)
		if err == sql.ErrNoRows{
			return	c.Status(fiber.StatusOK).JSON(fiber.Map{
				"message":"Data Trash tidak ditemukan",
				"success": true,
				"isTrash":	false,
			})
		}

		// Buat response pakai model 
    response := models.TrashResponse{ 
        Data: trash,
        Meta: models.MetaInfo{ 
            Page:   page, 
            Limit:  limit, 
            SortBy: sortBy, 
            Order:  order, 
            Search: search, 
        }, 
    } 
	
	return c.JSON(response)
}

func UpdateTrashService(c *fiber.Ctx, db *sql.DB, id string) error {
	username := c.Locals("username").(string) 
    pekerjaanID, err := strconv.Atoi(id) 
    if err != nil { 
        return c.Status(400).JSON(fiber.Map{ 
            "error": "ID tidak valid", 
        }) 
    } 
 
    log.Printf("Admin %s mengupdate pekerjaan ID %d", username, pekerjaanID)

	var trash models.UpdateTrash

	UpdatedTrash, err := repository.UpdateTrash(db, trash, id)
    
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"pesan": fmt.Sprintf("Data trash dengan ID %s tidak ditemukan untuk diperbarui.", id),
			})
		}
        
		fmt.Printf("Error saat memperbarui data trash ID %s: %v\n", id, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"pesan": "Error server: Gagal memperbarui data",
            // "detail_error": err.Error(), 
		})
	}

    // 4. Jika Sukses (data ditemukan dan diperbarui)
	return c.Status(fiber.StatusOK).JSON(UpdatedTrash)
}

func DeleteTrashService(c *fiber.Ctx, db *sql.DB, id string) error {
	username := c.Locals("username").(string) 
    trashID, err := strconv.Atoi(id) 
    if err != nil { 
        return c.Status(400).JSON(fiber.Map{ 
            "error": "ID tidak valid", 
        }) 
    } 

	log.Printf("Admin %s menghapus sampah ID %d", username, trashID)

	var trash models.Trash

	deletedTrash, err := repository.DeleteTrash(db, trash , id)
	if err == sql.ErrNoRows {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"pesan": "data sampah dengan ID " + id + " tidak ditemukan",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "pesan": "Data sampah berhasil dihapus",
        "data_dihapus": deletedTrash,
    })
}