package handlers

import (
	"database/sql"
	"net/http"
	"swiflet-backend/internal/config"
	"swiflet-backend/internal/database"
	"swiflet-backend/internal/models"
	"swiflet-backend/pkg/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	db       *database.DB
	config   *config.Config
	validate *validator.Validate
}

func NewAuthHandler(db *database.DB, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		db:       db,
		config:   cfg,
		validate: validator.New(),
	}
}

// Helper function to get int value or default
func getIntValue(ptr *int, defaultValue int) int {
	if ptr != nil {
		return *ptr
	}
	return defaultValue
}

// Register handles user registration
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid request body",
		})
		return
	}

	// Validate request
	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Validation failed",
		})
		return
	}

	// Check if email already exists
	var count int
	err := h.db.PostgreSQL.QueryRow("SELECT COUNT(*) FROM users WHERE email = $1", req.Email).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Database error",
		})
		return
	}

	if count > 0 {
		c.JSON(http.StatusConflict, models.ErrorResponse{
			Error: "Email already exists",
		})
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to hash password",
		})
		return
	}

	// Insert user
	var user models.User
	err = h.db.PostgreSQL.QueryRow(`
		INSERT INTO users (email, name, location, no_telp, password, img_profile, status, role, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, email, name, location, no_telp, img_profile, status, role, created_at
	`, req.Email, req.Name, req.Location, req.NoTelp, hashedPassword, req.ImgProfile, 
	   getIntValue(req.Status, 0), getIntValue(req.Role, 0), time.Now()).Scan(
		&user.ID, &user.Email, &user.Name, &user.Location, &user.NoTelp, 
		&user.ImgProfile, &user.Status, &user.Role, &user.CreatedAt,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to create user",
		})
		return
	}

	// Generate JWT token
	tokenString, _, err := utils.GenerateJWT(user.ID, user.Email, h.config.JWT.Secret, h.config.JWT.Expiry)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to generate token",
		})
		return
	}

	response := models.RegisterResponse{
		User:  user,
		Token: models.AuthToken(tokenString),
	}

	c.JSON(http.StatusCreated, response)
}

// Login handles user authentication
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid request body",
		})
		return
	}

	// Validate request
	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Validation failed",
		})
		return
	}

	// Get user by email
	var user models.User
	err := h.db.PostgreSQL.QueryRow(`
		SELECT id, email, name, location, no_telp, password, img_profile, status, role, created_at
		FROM users WHERE email = $1
	`, req.Email).Scan(
		&user.ID, &user.Email, &user.Name, &user.Location, &user.NoTelp, &user.Password, 
		&user.ImgProfile, &user.Status, &user.Role, &user.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Error: "Invalid credentials",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Database error",
		})
		return
	}

	// Check password
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: "Invalid credentials",
		})
		return
	}

	// Generate JWT token
	tokenString, _, err := utils.GenerateJWT(user.ID, user.Email, h.config.JWT.Secret, h.config.JWT.Expiry)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to generate token",
		})
		return
	}

	response := models.LoginResponse{
		Token: models.AuthToken(tokenString),
	}

	c.JSON(http.StatusOK, response)
}