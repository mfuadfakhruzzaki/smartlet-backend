package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"swiflet-backend/internal/database"
	"swiflet-backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	db       *database.DB
	validate *validator.Validate
}

func NewUserHandler(db *database.DB) *UserHandler {
	return &UserHandler{
		db:       db,
		validate: validator.New(),
	}
}

// ListUsers returns paginated list of users
func (h *UserHandler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage := 10
	offset := (page - 1) * perPage

	// Get total count
	var total int
	err := h.db.PostgreSQL.QueryRow("SELECT COUNT(*) FROM users").Scan(&total)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Database error",
		})
		return
	}

	// Get users
	rows, err := h.db.PostgreSQL.Query(`
		SELECT id, email, name, location, no_telp, img_profile, status, role, created_at
		FROM users
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`, perPage, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Database error",
		})
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID, &user.Email, &user.Name, &user.Location, &user.NoTelp,
			&user.ImgProfile, &user.Status, &user.Role, &user.CreatedAt,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error: "Database error",
			})
			return
		}
		users = append(users, user)
	}

	totalPages := (total + perPage - 1) / perPage
	response := models.PaginatedResponse[models.User]{
		Data:       users,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	}

	c.JSON(http.StatusOK, response)
}

// GetUser returns user by ID
func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid user ID",
		})
		return
	}

	var user models.User
	err = h.db.PostgreSQL.QueryRow(`
		SELECT id, email, name, location, no_telp, img_profile, status, role, created_at
		FROM users WHERE id = $1
	`, id).Scan(
		&user.ID, &user.Email, &user.Name, &user.Location, &user.NoTelp,
		&user.ImgProfile, &user.Status, &user.Role, &user.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error: "User not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Database error",
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUser updates user information
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid user ID",
		})
		return
	}

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid request body",
		})
		return
	}

	// Check if user exists
	var count int
	err = h.db.PostgreSQL.QueryRow("SELECT COUNT(*) FROM users WHERE id = $1", id).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Database error",
		})
		return
	}

	if count == 0 {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error: "User not found",
		})
		return
	}

	// Update user
	_, err = h.db.PostgreSQL.Exec(`
		UPDATE users 
		SET name = $1, location = $2, no_telp = $3, img_profile = $4, status = $5, role = $6
		WHERE id = $7
	`, user.Name, user.Location, user.NoTelp, user.ImgProfile, user.Status, user.Role, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to update user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

// DeleteUser deletes a user
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid user ID",
		})
		return
	}

	// Check if user exists
	var count int
	err = h.db.PostgreSQL.QueryRow("SELECT COUNT(*) FROM users WHERE id = $1", id).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Database error",
		})
		return
	}

	if count == 0 {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error: "User not found",
		})
		return
	}

	// Delete user
	_, err = h.db.PostgreSQL.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to delete user",
		})
		return
	}

	c.Status(http.StatusNoContent)
}