package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"
	"swiflet-backend/internal/database"
	"swiflet-backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type TagHandler struct {
	db       *database.DB
	validate *validator.Validate
}

func NewTagHandler(db *database.DB) *TagHandler {
	return &TagHandler{
		db:       db,
		validate: validator.New(),
	}
}

// ListTags returns paginated list of tags
func (h *TagHandler) ListTags(c *gin.Context) {
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
	countQuery := "SELECT COUNT(*) FROM tags"
	dataQuery := `
		SELECT id, name 
		FROM tags 
		ORDER BY name ASC 
		LIMIT $1 OFFSET $2
	`
	args := []interface{}{perPage, offset}

	if search != "" {
		countQuery += " WHERE name ILIKE $1"
		dataQuery = `
			SELECT id, name 
			FROM tags 
			WHERE name ILIKE $3
			ORDER BY name ASC 
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
			Error: "Failed to count tags",
		})
		return
	}

	// Get tags
	rows, err := h.db.PostgreSQL.Query(dataQuery, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to fetch tags",
		})
		return
	}
	defer rows.Close()

	var tags []models.Tags
	for rows.Next() {
		var tag models.Tags
		err := rows.Scan(&tag.ID, &tag.Name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error: "Failed to scan tag data",
			})
			return
		}
		tags = append(tags, tag)
	}

	// Handle empty results
	if tags == nil {
		tags = []models.Tags{}
	}

	totalPages := (total + perPage - 1) / perPage
	response := models.PaginatedResponse[models.Tags]{
		Data:       tags,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	}

	c.JSON(http.StatusOK, response)
}

// CreateTag creates a new tag
func (h *TagHandler) CreateTag(c *gin.Context) {
	var tag models.Tags
	if err := c.ShouldBindJSON(&tag); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid request body",
		})
		return
	}

	// Validate request
	if err := h.validate.Struct(tag); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Validation failed: " + err.Error(),
		})
		return
	}

	// Sanitize and validate tag name
	tag.Name = strings.TrimSpace(tag.Name)
	if tag.Name == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Tag name cannot be empty",
		})
		return
	}

	if len(tag.Name) < 2 || len(tag.Name) > 50 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Tag name must be between 2 and 50 characters",
		})
		return
	}

	// Check if tag already exists (case-insensitive)
	var count int
	err := h.db.PostgreSQL.QueryRow("SELECT COUNT(*) FROM tags WHERE LOWER(name) = LOWER($1)", tag.Name).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Database error",
		})
		return
	}

	if count > 0 {
		c.JSON(http.StatusConflict, models.ErrorResponse{
			Error: "Tag with this name already exists",
		})
		return
	}

	// Insert tag
	var newTag models.Tags
	err = h.db.PostgreSQL.QueryRow(`
		INSERT INTO tags (name)
		VALUES ($1)
		RETURNING id, name
	`, tag.Name).Scan(&newTag.ID, &newTag.Name)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to create tag",
		})
		return
	}

	c.JSON(http.StatusCreated, newTag)
}

// GetTag returns tag by ID
func (h *TagHandler) GetTag(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid tag ID",
		})
		return
	}

	var tag models.Tags
	err = h.db.PostgreSQL.QueryRow(`
		SELECT id, name
		FROM tags WHERE id = $1
	`, id).Scan(&tag.ID, &tag.Name)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error: "Tag not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Database error",
		})
		return
	}

	c.JSON(http.StatusOK, tag)
}

// UpdateTag updates a tag
func (h *TagHandler) UpdateTag(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid tag ID",
		})
		return
	}

	var tag models.Tags
	if err := c.ShouldBindJSON(&tag); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid request body",
		})
		return
	}

	// Validate request
	if err := h.validate.Struct(tag); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Validation failed: " + err.Error(),
		})
		return
	}

	// Sanitize and validate tag name
	tag.Name = strings.TrimSpace(tag.Name)
	if tag.Name == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Tag name cannot be empty",
		})
		return
	}

	if len(tag.Name) < 2 || len(tag.Name) > 50 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Tag name must be between 2 and 50 characters",
		})
		return
	}

	// Check if tag exists
	var currentName string
	err = h.db.PostgreSQL.QueryRow("SELECT name FROM tags WHERE id = $1", id).Scan(&currentName)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error: "Tag not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Database error",
		})
		return
	}

	// Check if new name conflicts with existing tag (excluding current tag)
	var count int
	err = h.db.PostgreSQL.QueryRow(
		"SELECT COUNT(*) FROM tags WHERE LOWER(name) = LOWER($1) AND id != $2", 
		tag.Name, id,
	).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Database error",
		})
		return
	}

	if count > 0 {
		c.JSON(http.StatusConflict, models.ErrorResponse{
			Error: "Tag with this name already exists",
		})
		return
	}

	// Update tag
	var updatedTag models.Tags
	err = h.db.PostgreSQL.QueryRow(`
		UPDATE tags 
		SET name = $1
		WHERE id = $2
		RETURNING id, name
	`, tag.Name, id).Scan(&updatedTag.ID, &updatedTag.Name)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to update tag",
		})
		return
	}

	c.JSON(http.StatusOK, updatedTag)
}

// DeleteTag deletes a tag
func (h *TagHandler) DeleteTag(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid tag ID",
		})
		return
	}

	// Check if tag exists
	var count int
	err = h.db.PostgreSQL.QueryRow("SELECT COUNT(*) FROM tags WHERE id = $1", id).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Database error",
		})
		return
	}

	if count == 0 {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error: "Tag not found",
		})
		return
	}

	// Check if tag is used by articles
	var articleCount int
	err = h.db.PostgreSQL.QueryRow("SELECT COUNT(*) FROM articles WHERE tag_id = $1", id).Scan(&articleCount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Database error",
		})
		return
	}

	if articleCount > 0 {
		c.JSON(http.StatusConflict, models.ErrorResponse{
			Error: "Cannot delete tag that is used by articles",
		})
		return
	}

	// Delete tag
	_, err = h.db.PostgreSQL.Exec("DELETE FROM tags WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to delete tag",
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetTagArticles returns articles associated with a tag
func (h *TagHandler) GetTagArticles(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid tag ID",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))

	// Validate pagination
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	offset := (page - 1) * perPage

	// Check if tag exists
	var tagName string
	err = h.db.PostgreSQL.QueryRow("SELECT name FROM tags WHERE id = $1", id).Scan(&tagName)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error: "Tag not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Database error",
		})
		return
	}

	// Get total count
	var total int
	err = h.db.PostgreSQL.QueryRow("SELECT COUNT(*) FROM articles WHERE tag_id = $1", id).Scan(&total)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to count articles",
		})
		return
	}

	// Get articles
	rows, err := h.db.PostgreSQL.Query(`
		SELECT id, title, content, cover_image, tag_id, status, created_at, updated_at
		FROM articles
		WHERE tag_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`, id, perPage, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to fetch articles",
		})
		return
	}
	defer rows.Close()

	var articles []models.Article
	for rows.Next() {
		var article models.Article
		err := rows.Scan(
			&article.ID, &article.Title, &article.Content, &article.CoverImage,
			&article.TagID, &article.Status, &article.CreatedAt, &article.UpdatedAt,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error: "Failed to scan article data",
			})
			return
		}
		articles = append(articles, article)
	}

	// Handle empty results
	if articles == nil {
		articles = []models.Article{}
	}

	totalPages := (total + perPage - 1) / perPage
	response := gin.H{
		"tag": gin.H{
			"id":   id,
			"name": tagName,
		},
		"articles": models.PaginatedResponse[models.Article]{
			Data:       articles,
			Page:       page,
			PerPage:    perPage,
			Total:      total,
			TotalPages: totalPages,
		},
	}

	c.JSON(http.StatusOK, response)
}