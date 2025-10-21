package routes

import (
	"database/sql"
	"latihan_uts_2/app/models"
	"latihan_uts_2/app/services"
	"latihan_uts_2/middleware"
	"latihan_uts_2/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRoutes(app *fiber.App, postgresDB *sql.DB, mongoDB *mongo.Client) {
	api := app.Group("/api")
	 
    // Public routes (tidak perlu login)
	api.Group("") 
	api.Post("/login", func(c *fiber.Ctx) error {
		return login(c, postgresDB) 
	})

	// Protected routes (perlu login) 
	protected := api.Group("", middleware.AuthRequired()) 
    protected.Get("/profile", getProfile) 

	// --- Grup Rute Alumni ---
	alumniGroup := protected.Group("/alumni")
	// GET /alumni
	alumniGroup.Get("/", func(c *fiber.Ctx) error {
		return services.GetAllAlumniService(c, postgresDB)
	})
	// GET /alumni/:id
	alumniGroup.Get("/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		return services.GetAlumniByIDService(c, postgresDB, id)
	})
	// POST /alumni
	alumniGroup.Post("/", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return services.CreateAlumniService(c, postgresDB)
	})
	// PUT /alumni/:id
	alumniGroup.Put("/:id", middleware.AdminOnly(),func(c *fiber.Ctx) error {
		id := c.Params("id")
		return services.UpdateAlumniService(c, postgresDB, id)
	})
	// DELETE /alumni/:id
	alumniGroup.Delete("/:id", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		id := c.Params("id")
		return services.SoftDeleteAlumniService(c, postgresDB, id)
	})

	// --- Grup Rute Pekerjaan ---
	pekerjaanGroup := protected.Group("/pekerjaan") 
	// GET /pekerjaan
	pekerjaanGroup.Get("/", func(c *fiber.Ctx) error {
		return services.GetAllPekerjaanService(c, postgresDB)
	})
	// GET /pekerjaan/:id
	pekerjaanGroup.Get("/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		return services.GetPekerjaanByIDService(c, postgresDB, id)
	})
	// POST /pekerjaan
	pekerjaanGroup.Post("/", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return services.CreatePekerjaanService(c, postgresDB)
	})
	// PUT /pekerjaan/:id
	pekerjaanGroup.Put("/:id", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		id := c.Params("id")
		return services.UpdatePekerjaanService(c, postgresDB, id)
	})
	// DELETE /pekerjaan/:id (SoftDelete)
	pekerjaanGroup.Delete("/:id", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		id := c.Params("id")
		return services.SoftDeletePekerjaanService(c, postgresDB, id)
	})

	// --- Trash ---
	trashed := protected.Group("/trash")
    trashed.Get("/", func(c *fiber.Ctx) error {
		return services.GetAllTrashService(c, postgresDB)
	})

	trashed.Put("/:id", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		id := c.Params("id")
		return services.UpdateTrashService(c, postgresDB, id)
	})

	trashed.Delete("/:id", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		id := c.Params("id")
		return services.DeleteTrashService(c, postgresDB, id)
	})
}

func login(c *fiber.Ctx, postgresDB *sql.DB) error {
	 var req models.LoginRequest 
	 if err := c.BodyParser(&req); err != nil { 
        return c.Status(400).JSON(fiber.Map{ 
            "error": "Request body tidak valid", 
        }) 
    } 

	// Validasi input 
    if req.Username == "" || req.Password == "" { 
        return c.Status(400).JSON(fiber.Map{ 
            "error": "Username dan password harus diisi", 
        }) 
    }

	   // Cari user di database 
    var user models.User 
    var passwordHash string 
    err := postgresDB.QueryRow(` 
        SELECT id, username, email, password_hash, role, created_at 
        FROM users  
        WHERE username = $1 OR email = $1 
    `, req.Username).Scan( 
        &user.ID, &user.Username, &user.Email, &passwordHash, &user.Role, &user.CreatedAt,
	)

	if err != nil { 
	if err == sql.ErrNoRows { 
		return c.Status(401).JSON(fiber.Map{ 
			"error": "Username atau password salah", 
		}) 
	} 
	return c.Status(500).JSON(fiber.Map{ 
			"error": "Error database", 
		}) 
    }

    // Check password 
    if !utils.CheckPassword(req.Password, passwordHash) { 
        return c.Status(401).JSON(fiber.Map{ 
            "error": "Username atau password salah", 
        }) 
    } 

	// Generate JWT token 
    token, err := utils.GenerateToken(user) 
    if err != nil { 
        return c.Status(500).JSON(fiber.Map{ 
            "error": "Gagal generate token", 
        }) 
    } 
 
    response := models.LoginResponse{ 
        User:  user, 
        Token: token, 
    } 
 
    return c.JSON(fiber.Map{ 
        "success": true, 
        "message": "Login berhasil", 
        "data":    response, 
    }) 
}

func getProfile(c *fiber.Ctx) error { 
    userID := c.Locals("user_id").(int) 
    username := c.Locals("username").(string) 
    role := c.Locals("role").(string) 
 
    return c.JSON(fiber.Map{ 
        "success": true, 
        "message": "Profile berhasil diambil", 
        "data": fiber.Map{ 
            "user_id":  userID, 
            "username": username, 
            "role":     role, 
        }, 
    }) 
}