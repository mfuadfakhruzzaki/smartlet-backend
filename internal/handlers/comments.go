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

type CommentHandler struct {
	db       *database.DB
	validate *validator.Validate
}

func NewCommentHandler(db *database.DB) *CommentHandler {
	return &CommentHandler{
		db:       db,
		validate: validator.New(),
	}
}

// ListComments returns paginated list of comments for an article
func (h *CommentHandler) ListComments(c *gin.Context) {
	articleID, err := strconv.Atoi(c.Param("article_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid article ID",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))

	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	offset := (page - 1) * perPage

	// Check if article exists
	var articleExists int
	err = h.db.PostgreSQL.QueryRow("SELECT COUNT(*) FROM articles WHERE id = $1", articleID).Scan(&articleExists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Database error",
		})
		return
	}

	if articleExists == 0 {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error: "Article not found",
		})
		return
	}

	// Get total count
	var total int
	err = h.db.PostgreSQL.QueryRow("SELECT COUNT(*) FROM comments WHERE article_id = $1", articleID).Scan(&total)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to count comments",
		})
		return
	}

	// Get comments with user information
	rows, err := h.db.PostgreSQL.Query(`
		SELECT c.id, c.article_id, c.user_id, c.content, c.created_at, c.updated_at,
		       u.name as user_name, u.email as user_email, u.img_profile
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.article_id = $1
		ORDER BY c.created_at DESC
		LIMIT $2 OFFSET $3
	`, articleID, perPage, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to fetch comments",
		})
		return
	}
	defer rows.Close()

	type CommentWithUser struct {
		models.Comment
		UserName     string  `json:"user_name"`
		UserEmail    string  `json:"user_email"`
		UserImage    *string `json:"user_image"`
	}

	var comments []CommentWithUser
	for rows.Next() {
		var comment CommentWithUser
		err := rows.Scan(
			&comment.ID, &comment.ArticleID, &comment.UserID, &comment.Content, 
			&comment.CreatedAt, &comment.UpdatedAt, &comment.UserName, 
			&comment.UserEmail, &comment.UserImage,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error: "Failed to scan comment data",
			})
			return
		}
		comments = append(comments, comment)
	}

	// Handle empty results
	if comments == nil {
		comments = []CommentWithUser{}
	}

	totalPages := (total + perPage - 1) / perPage
	response := models.PaginatedResponse[CommentWithUser]{
		Data:       comments,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	}

	c.JSON(http.StatusOK, response)
}

// CreateComment creates a new comment
func (h *CommentHandler) CreateComment(c *gin.Context) {
	articleID, err := strconv.Atoi(c.Param("article_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid article ID",
		})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: "User not authenticated",
		})
		return
	}

	var request struct {
		Content string `json:"content" validate:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid request body",
		})
		return
	}

	// Validate request
	if err := h.validate.Struct(request); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Validation failed: " + err.Error(),
		})
		return
	}

	// Sanitize content
	content := strings.TrimSpace(request.Content)
	if content == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Comment content cannot be empty",
		})
		return
	}

	if len(content) < 1 || len(content) > 1000 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Comment content must be between 1 and 1000 characters",
		})
		return
	}

	// Check if article exists
	var articleExists int
	err = h.db.PostgreSQL.QueryRow("SELECT COUNT(*) FROM articles WHERE id = $1", articleID).Scan(&articleExists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Database error",
		})
		return
	}

	if articleExists == 0 {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error: "Article not found",
		})
		return
	}

	// Insert comment
	var comment models.Comment
	now := time.Now()
	err = h.db.PostgreSQL.QueryRow(`
		INSERT INTO comments (article_id, user_id, content, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, article_id, user_id, content, created_at, updated_at
	`, articleID, userID, content, now, now).Scan(
		&comment.ID, &comment.ArticleID, &comment.UserID, 
		&comment.Content, &comment.CreatedAt, &comment.UpdatedAt,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to create comment",
		})
		return
	}

	c.JSON(http.StatusCreated, comment)
}

// GetComment returns a specific comment
func (h *CommentHandler) GetComment(c *gin.Context) {
	articleID, err := strconv.Atoi(c.Param("article_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid article ID",
		})
		return
	}

	commentID, err := strconv.Atoi(c.Param("comment_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid comment ID",
		})
		return
	}

	type CommentWithUser struct {
		models.Comment
		UserName     string  `json:"user_name"`
		UserEmail    string  `json:"user_email"`
		UserImage    *string `json:"user_image"`
	}

	var comment CommentWithUser
	err = h.db.PostgreSQL.QueryRow(`
		SELECT c.id, c.article_id, c.user_id, c.content, c.created_at, c.updated_at,
		       u.name as user_name, u.email as user_email, u.img_profile
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.id = $1 AND c.article_id = $2
	`, commentID, articleID).Scan(
		&comment.ID, &comment.ArticleID, &comment.UserID, &comment.Content,
		&comment.CreatedAt, &comment.UpdatedAt, &comment.UserName,
		&comment.UserEmail, &comment.UserImage,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error: "Comment not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Database error",
		})
		return
	}

	c.JSON(http.StatusOK, comment)
}

// UpdateComment updates a comment
func (h *CommentHandler) UpdateComment(c *gin.Context) {
	articleID, err := strconv.Atoi(c.Param("article_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid article ID",
		})
		return
	}

	commentID, err := strconv.Atoi(c.Param("comment_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid comment ID",
		})
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: "User not authenticated",
		})
		return
	}

	var request struct {
		Content string `json:"content" validate:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid request body",
		})
		return
	}

	// Validate request
	if err := h.validate.Struct(request); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Validation failed: " + err.Error(),
		})
		return
	}

	// Sanitize content
	content := strings.TrimSpace(request.Content)
	if content == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Comment content cannot be empty",
		})
		return
	}

	if len(content) < 1 || len(content) > 1000 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Comment content must be between 1 and 1000 characters",
		})
		return
	}

	// Check if comment exists and belongs to user
	var existingUserID int
	err = h.db.PostgreSQL.QueryRow(`
		SELECT user_id FROM comments 
		WHERE id = $1 AND article_id = $2
	`, commentID, articleID).Scan(&existingUserID)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error: "Comment not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Database error",
		})
		return
	}

	// Check if user owns the comment
	if existingUserID != userID.(int) {
		c.JSON(http.StatusForbidden, models.ErrorResponse{
			Error: "You can only edit your own comments",
		})
		return
	}

	// Update comment
	var updatedComment models.Comment
	err = h.db.PostgreSQL.QueryRow(`
		UPDATE comments 
		SET content = $1, updated_at = $2
		WHERE id = $3 AND article_id = $4
		RETURNING id, article_id, user_id, content, created_at, updated_at
	`, content, time.Now(), commentID, articleID).Scan(
		&updatedComment.ID, &updatedComment.ArticleID, &updatedComment.UserID,
		&updatedComment.Content, &updatedComment.CreatedAt, &updatedComment.UpdatedAt,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to update comment",
		})
		return
	}

	c.JSON(http.StatusOK, updatedComment)
}

// DeleteComment deletes a comment
func (h *CommentHandler) DeleteComment(c *gin.Context) {
	articleID, err := strconv.Atoi(c.Param("article_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid article ID",
		})
		return
	}

	commentID, err := strconv.Atoi(c.Param("comment_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid comment ID",
		})
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: "User not authenticated",
		})
		return
	}

	// Check if comment exists and belongs to user
	var existingUserID int
	err = h.db.PostgreSQL.QueryRow(`
		SELECT user_id FROM comments 
		WHERE id = $1 AND article_id = $2
	`, commentID, articleID).Scan(&existingUserID)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error: "Comment not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Database error",
		})
		return
	}

	// Check if user owns the comment (or could be admin in future)
	if existingUserID != userID.(int) {
		c.JSON(http.StatusForbidden, models.ErrorResponse{
			Error: "You can only delete your own comments",
		})
		return
	}

	// Delete comment
	_, err = h.db.PostgreSQL.Exec(`
		DELETE FROM comments 
		WHERE id = $1 AND article_id = $2
	`, commentID, articleID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to delete comment",
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetUserComments returns all comments by a specific user
func (h *CommentHandler) GetUserComments(c *gin.Context) {
	targetUserID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid user ID",
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

	// Check if user exists
	var userExists int
	err = h.db.PostgreSQL.QueryRow("SELECT COUNT(*) FROM users WHERE id = $1", targetUserID).Scan(&userExists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Database error",
		})
		return
	}

	if userExists == 0 {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error: "User not found",
		})
		return
	}

	// Get total count
	var total int
	err = h.db.PostgreSQL.QueryRow("SELECT COUNT(*) FROM comments WHERE user_id = $1", targetUserID).Scan(&total)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to count comments",
		})
		return
	}

	// Get comments with article information
	rows, err := h.db.PostgreSQL.Query(`
		SELECT c.id, c.article_id, c.user_id, c.content, c.created_at, c.updated_at,
		       a.title as article_title
		FROM comments c
		JOIN articles a ON c.article_id = a.id
		WHERE c.user_id = $1
		ORDER BY c.created_at DESC
		LIMIT $2 OFFSET $3
	`, targetUserID, perPage, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to fetch comments",
		})
		return
	}
	defer rows.Close()

	type CommentWithArticle struct {
		models.Comment
		ArticleTitle string `json:"article_title"`
	}

	var comments []CommentWithArticle
	for rows.Next() {
		var comment CommentWithArticle
		err := rows.Scan(
			&comment.ID, &comment.ArticleID, &comment.UserID, &comment.Content,
			&comment.CreatedAt, &comment.UpdatedAt, &comment.ArticleTitle,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error: "Failed to scan comment data",
			})
			return
		}
		comments = append(comments, comment)
	}

	// Handle empty results
	if comments == nil {
		comments = []CommentWithArticle{}
	}

	totalPages := (total + perPage - 1) / perPage
	response := models.PaginatedResponse[CommentWithArticle]{
		Data:       comments,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	}

	c.JSON(http.StatusOK, response)
}