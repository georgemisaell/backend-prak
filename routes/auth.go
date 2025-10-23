package routes

import (
	"database/sql"
	"latihan_uts_2/app/models"
	"latihan_uts_2/utils"

	"github.com/gofiber/fiber/v2"
)

func Login(c *fiber.Ctx, postgreSQL *sql.DB) error {
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
	err := postgreSQL.QueryRow(` 
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

func GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	username := c.Locals("username").(string)
	role := c.Locals("role").(string)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Profile berhasil diambil",
		"data": fiber.Map{
			"userfiber": userID,
			"username": username,
			"role": role,
		},
	})
}