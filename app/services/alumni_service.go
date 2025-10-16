package services

import (
	"database/sql"
	"latihan_uts_2/app/models"
	"latihan_uts_2/app/repository"
	"log"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func GetAllAlumniService(c *fiber.Ctx, db *sql.DB) error {
	// Log siapa yang mengakses 
    username := c.Locals("username").(string) 
    log.Printf("User %s mengakses GET /api/alumni", username) 

	// Pagination
	page, _ := strconv.Atoi(c.Query("page", "1")) 
    limit, _ := strconv.Atoi(c.Query("limit", "10")) 
    sortBy := c.Query("sortBy", "id") 
    order := c.Query("order", "asc") 
    search := c.Query("search", "") 
	offset := (page - 1) * limit

	// Validasi input 
    sortByWhitelist := map[string]bool{"id": true, "nim": true, "nama": true, "created_at": true} 
    if !sortByWhitelist[sortBy] { 
        sortBy = "id" 
    } 
    if strings.ToLower(order) != "desc" { 
        order = "asc" 
    }

	// Ambil data dari repository 
    // alumni, err := repository.GetAllAlumni(search, sortBy, order, limit, offset) 
    // if err != nil { 
    //     return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch alumni"}) 
    // } 
 
    total, err := repository.CountAlumniRepo(search, db) 
    if err != nil { 
        return c.Status(500).JSON(fiber.Map{"error": "Failed to count alumni"}) 
    }

	alumniList, err := repository.GetAllAlumni(search, sortBy, order, limit, offset, db)
	if err != nil {
		if err == sql.ErrNoRows{
			return	c.Status(fiber.StatusOK).JSON(fiber.Map{
				"message":"Alumni tidak ditemukan",
				"success": true,
				"isAlumni":	false,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal cek alumni karena" + err.Error(),
			"Success": false,
		})
	}

	// Buat response pakai model 
    response := models.AlumniResponse{ 
        Data: alumniList,
        Meta: models.MetaInfo{ 
            Page:   page, 
            Limit:  limit, 
            Total:  total, 
            Pages:  (total + limit - 1) / limit, 
            SortBy: sortBy, 
            Order:  order, 
            Search: search, 
        }, 
    } 
	
	return c.JSON(response)
}

func GetAlumniByIDService(c *fiber.Ctx, db *sql.DB, id string) error{
   	username := c.Locals("username").(string) 
    alumniID, err := strconv.Atoi(id)
    if err != nil { 
        return c.Status(400).JSON(fiber.Map{ 
            "error": "ID tidak valid", 
        }) 
    } 
 
    log.Printf("User %s mengakses GET /api/alumni/%d", username, alumniID) 

	alumni, err := repository.GetAlumniByID(db, id)
	if err != nil {
		if err == sql.ErrNoRows{
			return	c.Status(fiber.StatusOK).JSON(fiber.Map{
				"message":"ID alumni tidak ditemukan",	
				"success": true,	
				"isAlumni":	false,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal cek alumni karena" + err.Error(),
			"Success": false,
		})
	}
	
	return c.JSON(fiber.Map{
		"success": true,
		"data": alumni,
		"message":"Data alumni berhasil diambil",
	})
}

func CreateAlumniService(c *fiber.Ctx, db *sql.DB) error{
	username := c.Locals("username").(string) 
    log.Printf("Admin %s menambah alumni baru", username) 

	var alumni models.CreateAlumni

    // Ambil JSON body dari request
    if err := c.BodyParser(&alumni); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "success": false,
            "message": "Request body tidak valid: " + err.Error(),
        })
    }

    // Simpan ke DB lewat repository
    createdAlumni, err := repository.CreateAlumni(db, alumni)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "success": false,
            "message": "Gagal menambahkan alumni: " + err.Error(),
        })
    }
	
    return c.Status(fiber.StatusCreated).JSON(fiber.Map{
        "success": true,
        "message": "Data alumni berhasil ditambahkan",
        "data":    createdAlumni,
    })
}

func UpdateAlumniService(c *fiber.Ctx, db *sql.DB, id string) error{
	username := c.Locals("username").(string) 
    alumniID, err := strconv.Atoi(id) 
    if err != nil { 
        return c.Status(400).JSON(fiber.Map{ 
            "error": "ID tidak valid", 
        }) 
    } 
 
    log.Printf("Admin %s mengupdate alumni ID %d", username, alumniID)

	var alumni models.UpdateAlumni

    // Ambil JSON body dari request
    if err := c.BodyParser(&alumni); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "success": false,
            "message": "Request body tidak valid: " + err.Error(),
        })
    }

    // Simpan ke DB lewat repository
    updatedAlumni, err := repository.UpdateAlumni(db, alumni, id)
    if err != nil {
        if err == sql.ErrNoRows {
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
                "success": false,
                "message": "Alumni dengan ID " + id + " tidak ditemukan",
            })
        }
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "success": false,
            "message": "Gagal update alumni: " + err.Error(),
        })
    }
	
    return c.Status(fiber.StatusCreated).JSON(fiber.Map{
        "success": true,
        "message": "Data alumni berhasil diupdate",
        "data":    updatedAlumni,
    })
}

func DeleteAlumniService(c *fiber.Ctx, db *sql.DB, id string) error{
	username := c.Locals("username").(string) 
    alumniID, err := strconv.Atoi(id) 
    if err != nil { 
        return c.Status(400).JSON(fiber.Map{ 
            "error": "ID tidak valid", 
        }) 
    } 

	log.Printf("Admin %s menghapus alumni ID %d", username, alumniID)

	alumni, err := repository.DeleteAlumni(db, id)
	if err != nil {
		if err == sql.ErrNoRows{
			return	c.Status(fiber.StatusOK).JSON(fiber.Map{
				"message":"ID alumni tidak ditemukan",	
				"success": true,	
				"isAlumni":	false,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal hapus alumni karena" + err.Error(),
			"Success": false,
		})
	}
	
	return c.JSON(fiber.Map{
		"success": true,
		"data": alumni,
		"message":"Data alumni berhasil dihapus",
	})
}