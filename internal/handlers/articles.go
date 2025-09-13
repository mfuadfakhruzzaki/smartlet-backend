package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"swiflet-backend/internal/database"
	"swiflet-backend/internal/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ArticleHandler struct {
	db       *database.DB
	validate *validator.Validate
}

func NewArticleHandler(db *database.DB) *ArticleHandler {
	return &ArticleHandler{
		db:       db,
		validate: validator.New(),
	}
}

// ListArticles returns paginated list of articles
func (h *ArticleHandler) ListArticles(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage := 10
	offset := (page - 1) * perPage

	// Get total count
	var total int
	err := h.db.PostgreSQL.QueryRow("SELECT COUNT(*) FROM articles").Scan(&total)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Database error",
		})
		return
	}

	// Get articles
	rows, err := h.db.PostgreSQL.Query(`
		SELECT id, title, content, cover_image, tag_id, status, created_at, updated_at
		FROM articles
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

	var articles []models.Article
	for rows.Next() {
		var article models.Article
		err := rows.Scan(
			&article.ID, &article.Title, &article.Content, &article.CoverImage,
			&article.TagID, &article.Status, &article.CreatedAt, &article.UpdatedAt,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error: "Database error",
			})
			return
		}
		articles = append(articles, article)
	}

	totalPages := (total + perPage - 1) / perPage
	response := models.PaginatedResponse[models.Article]{
		Data:       articles,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	}

	c.JSON(http.StatusOK, gin.H{"data": response.Data})
}

// CreateArticle creates a new article
func (h *ArticleHandler) CreateArticle(c *gin.Context) {
	var article models.Article
	if err := c.ShouldBindJSON(&article); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid request body",
		})
		return
	}

	// Validate request
	if err := h.validate.Struct(article); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Validation failed",
		})
		return
	}

	// Insert article
	_, err := h.db.PostgreSQL.Exec(`
		INSERT INTO articles (title, content, cover_image, tag_id, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, article.Title, article.Content, article.CoverImage, article.TagID, article.Status, time.Now(), time.Now())

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to create article",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Article created successfully"})
}

// GetArticle returns article by ID
func (h *ArticleHandler) GetArticle(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid article ID",
		})
		return
	}

	var article models.Article
	err = h.db.PostgreSQL.QueryRow(`
		SELECT id, title, content, cover_image, tag_id, status, created_at, updated_at
		FROM articles WHERE id = $1
	`, id).Scan(
		&article.ID, &article.Title, &article.Content, &article.CoverImage,
		&article.TagID, &article.Status, &article.CreatedAt, &article.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error: "Article not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Database error",
		})
		return
	}

	c.JSON(http.StatusOK, article)
}

// UpdateArticle updates an article
func (h *ArticleHandler) UpdateArticle(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid article ID",
		})
		return
	}

	var article models.Article
	if err := c.ShouldBindJSON(&article); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid request body",
		})
		return
	}

	// Check if article exists
	var count int
	err = h.db.PostgreSQL.QueryRow("SELECT COUNT(*) FROM articles WHERE id = $1", id).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Database error",
		})
		return
	}

	if count == 0 {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error: "Article not found",
		})
		return
	}

	// Update article
	_, err = h.db.PostgreSQL.Exec(`
		UPDATE articles 
		SET title = $1, content = $2, cover_image = $3, tag_id = $4, status = $5, updated_at = $6
		WHERE id = $7
	`, article.Title, article.Content, article.CoverImage, article.TagID, article.Status, time.Now(), id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to update article",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Article updated successfully"})
}

// DeleteArticle deletes an article
func (h *ArticleHandler) DeleteArticle(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid article ID",
		})
		return
	}

	// Check if article exists
	var count int
	err = h.db.PostgreSQL.QueryRow("SELECT COUNT(*) FROM articles WHERE id = $1", id).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Database error",
		})
		return
	}

	if count == 0 {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error: "Article not found",
		})
		return
	}

	// Delete article
	_, err = h.db.PostgreSQL.Exec("DELETE FROM articles WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to delete article",
		})
		return
	}

	c.Status(http.StatusNoContent)
}