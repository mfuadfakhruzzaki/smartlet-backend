package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"
	"swiflet-backend/internal/database"
	"swiflet-backend/internal/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type EBookHandler struct {
	db       *database.DB
	validate *validator.Validate
}

func NewEBookHandler(db *database.DB) *EBookHandler {
	return &EBookHandler{
		db:       db,
		validate: validator.New(),
	}
}

// ListEBooks returns paginated list of ebooks
func (h *EBookHandler) ListEBooks(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	search := strings.TrimSpace(c.DefaultQuery("search", ""))

	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	offset := (page - 1) * perPage

	// Build query with optional search
	countQuery := "SELECT COUNT(*) FROM ebooks"
	dataQuery := `
		SELECT id, title, file_path, thumbnail_path, created_at, updated_at
		FROM ebooks 
		ORDER BY created_at DESC 
		LIMIT $1 OFFSET $2
	`
	args := []interface{}{perPage, offset}

	if search != "" {
		countQuery += " WHERE title ILIKE $1"
		dataQuery = `
			SELECT id, title, file_path, thumbnail_path, created_at, updated_at
			FROM ebooks 
			WHERE title ILIKE $3
			ORDER BY created_at DESC 
			LIMIT $1 OFFSET $2
		`
		args = append(args, "%"+search+"%")
	}

	// Get total count
	var total int
	var err error
	if search != "" {
		err = h.db.PostgreSQL.QueryRow(countQuery, "%"+search+"%").Scan(&total)
	} else {
		err = h.db.PostgreSQL.QueryRow(countQuery).Scan(&total)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to count ebooks",
		})
		return
	}

	// Get ebooks
	rows, err := h.db.PostgreSQL.Query(dataQuery, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to fetch ebooks",
		})
		return
	}
	defer rows.Close()

	var ebooks []models.EBook
	for rows.Next() {
		var ebook models.EBook
		err := rows.Scan(
			&ebook.ID, &ebook.Title, &ebook.FilePath, &ebook.ThumbnailPath,
			&ebook.CreatedAt, &ebook.UpdatedAt,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error: "Failed to scan ebook data",
			})
			return
		}
		ebooks = append(ebooks, ebook)
	}

	// Handle empty results
	if ebooks == nil {
		ebooks = []models.EBook{}
	}

	totalPages := (total + perPage - 1) / perPage
	response := models.PaginatedResponse[models.EBook]{
		Data:       ebooks,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	}

	c.JSON(http.StatusOK, response)
}

// CreateEBook creates a new ebook
func (h *EBookHandler) CreateEBook(c *gin.Context) {
	var ebook models.EBook
	if err := c.ShouldBindJSON(&ebook); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid request body",
		})
		return
	}

	// Validate request
	if err := h.validate.Struct(ebook); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Validation failed: " + err.Error(),
		})
		return
	}

	// Sanitize and validate fields
	ebook.Title = strings.TrimSpace(ebook.Title)
	ebook.FilePath = strings.TrimSpace(ebook.FilePath)

	if ebook.Title == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "EBook title cannot be empty",
		})
		return
	}

	if len(ebook.Title) < 3 || len(ebook.Title) > 255 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "EBook title must be between 3 and 255 characters",
		})
		return
	}

	if ebook.FilePath == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "File path cannot be empty",
		})
		return
	}

	// Check if title already exists
	var count int
	err := h.db.PostgreSQL.QueryRow("SELECT COUNT(*) FROM ebooks WHERE LOWER(title) = LOWER($1)", ebook.Title).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Database error",
		})
		return
	}

	if count > 0 {
		c.JSON(http.StatusConflict, models.ErrorResponse{
			Error: "EBook with this title already exists",
		})
		return
	}

	// Insert ebook
	var newEBook models.EBook
	now := time.Now()
	err = h.db.PostgreSQL.QueryRow(`
		INSERT INTO ebooks (title, file_path, thumbnail_path, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, title, file_path, thumbnail_path, created_at, updated_at
	`, ebook.Title, ebook.FilePath, ebook.ThumbnailPath, now, now).Scan(
		&newEBook.ID, &newEBook.Title, &newEBook.FilePath, 
		&newEBook.ThumbnailPath, &newEBook.CreatedAt, &newEBook.UpdatedAt,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to create ebook",
		})
		return
	}

	c.JSON(http.StatusCreated, newEBook)
}

// GetEBook returns ebook by ID
func (h *EBookHandler) GetEBook(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid ebook ID",
		})
		return
	}

	var ebook models.EBook
	err = h.db.PostgreSQL.QueryRow(`
		SELECT id, title, file_path, thumbnail_path, created_at, updated_at
		FROM ebooks WHERE id = $1
	`, id).Scan(
		&ebook.ID, &ebook.Title, &ebook.FilePath, &ebook.ThumbnailPath,
		&ebook.CreatedAt, &ebook.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error: "EBook not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Database error",
		})
		return
	}

	c.JSON(http.StatusOK, ebook)
}

// UpdateEBook updates an ebook
func (h *EBookHandler) UpdateEBook(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid ebook ID",
		})
		return
	}

	var ebook models.EBook
	if err := c.ShouldBindJSON(&ebook); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid request body",
		})
		return
	}

	// Validate request
	if err := h.validate.Struct(ebook); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Validation failed: " + err.Error(),
		})
		return
	}

	// Sanitize and validate fields
	ebook.Title = strings.TrimSpace(ebook.Title)
	ebook.FilePath = strings.TrimSpace(ebook.FilePath)

	if ebook.Title == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "EBook title cannot be empty",
		})
		return
	}

	if len(ebook.Title) < 3 || len(ebook.Title) > 255 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "EBook title must be between 3 and 255 characters",
		})
		return
	}

	if ebook.FilePath == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "File path cannot be empty",
		})
		return
	}

	// Check if ebook exists
	var currentTitle string
	err = h.db.PostgreSQL.QueryRow("SELECT title FROM ebooks WHERE id = $1", id).Scan(&currentTitle)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error: "EBook not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Database error",
		})
		return
	}

	// Check if new title conflicts with existing ebook (excluding current ebook)
	var count int
	err = h.db.PostgreSQL.QueryRow(
		"SELECT COUNT(*) FROM ebooks WHERE LOWER(title) = LOWER($1) AND id != $2", 
		ebook.Title, id,
	).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Database error",
		})
		return
	}

	if count > 0 {
		c.JSON(http.StatusConflict, models.ErrorResponse{
			Error: "EBook with this title already exists",
		})
		return
	}

	// Update ebook
	var updatedEBook models.EBook
	err = h.db.PostgreSQL.QueryRow(`
		UPDATE ebooks 
		SET title = $1, file_path = $2, thumbnail_path = $3, updated_at = $4
		WHERE id = $5
		RETURNING id, title, file_path, thumbnail_path, created_at, updated_at
	`, ebook.Title, ebook.FilePath, ebook.ThumbnailPath, time.Now(), id).Scan(
		&updatedEBook.ID, &updatedEBook.Title, &updatedEBook.FilePath,
		&updatedEBook.ThumbnailPath, &updatedEBook.CreatedAt, &updatedEBook.UpdatedAt,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to update ebook",
		})
		return
	}

	c.JSON(http.StatusOK, updatedEBook)
}

// DeleteEBook deletes an ebook
func (h *EBookHandler) DeleteEBook(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid ebook ID",
		})
		return
	}

	// Check if ebook exists
	var filePath string
	var thumbnailPath *string
	err = h.db.PostgreSQL.QueryRow("SELECT file_path, thumbnail_path FROM ebooks WHERE id = $1", id).Scan(&filePath, &thumbnailPath)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error: "EBook not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Database error",
		})
		return
	}

	// Delete ebook record
	_, err = h.db.PostgreSQL.Exec("DELETE FROM ebooks WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to delete ebook",
		})
		return
	}

	// TODO: Delete actual files from storage when S3 service is implemented
	// For now, just return success

	c.Status(http.StatusNoContent)
}

// DownloadEBook handles ebook download
func (h *EBookHandler) DownloadEBook(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid ebook ID",
		})
		return
	}

	// Get ebook info
	var ebook models.EBook
	err = h.db.PostgreSQL.QueryRow(`
		SELECT id, title, file_path, thumbnail_path, created_at, updated_at
		FROM ebooks WHERE id = $1
	`, id).Scan(
		&ebook.ID, &ebook.Title, &ebook.FilePath, &ebook.ThumbnailPath,
		&ebook.CreatedAt, &ebook.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error: "EBook not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Database error",
		})
		return
	}

	// TODO: Implement actual file serving from S3 or local storage
	// For now, return the file path info
	c.JSON(http.StatusOK, gin.H{
		"ebook":        ebook,
		"download_url": ebook.FilePath, // This should be a signed URL in production
		"message":      "File download functionality will be implemented with S3 service",
	})
}