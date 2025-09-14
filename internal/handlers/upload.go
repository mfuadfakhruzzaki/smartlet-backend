package handlers

import (
	"net/http"
	"strconv"
	"swiflet-backend/internal/database"
	"swiflet-backend/internal/models"
	"swiflet-backend/internal/services"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	db        *database.DB
	s3Service *services.S3Service
}

func NewUploadHandler(db *database.DB, s3Service *services.S3Service) *UploadHandler {
	return &UploadHandler{
		db:        db,
		s3Service: s3Service,
	}
}

// UploadUserProfile uploads user profile image
func (h *UploadHandler) UploadUserProfile(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: "User not authenticated",
		})
		return
	}

	// Get uploaded file
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "No file uploaded",
		})
		return
	}
	defer file.Close()

	// Upload to S3
	result, err := h.s3Service.UploadUserProfileImage(file, header, userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to upload file: " + err.Error(),
		})
		return
	}

	// Update user profile image in database
	_, err = h.db.PostgreSQL.Exec(
		"UPDATE users SET img_profile = $1 WHERE id = $2",
		result.URL, userID,
	)
	if err != nil {
		// If database update fails, try to delete the uploaded file
		h.s3Service.DeleteFile(result.Key)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to update user profile",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile image uploaded successfully",
		"url":     result.URL,
		"size":    result.Size,
	})
}

// UploadArticleCover uploads article cover image
func (h *UploadHandler) UploadArticleCover(c *gin.Context) {
	articleID, err := strconv.Atoi(c.Param("article_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid article ID",
		})
		return
	}

	// Check if article exists
	var count int
	err = h.db.PostgreSQL.QueryRow("SELECT COUNT(*) FROM articles WHERE id = $1", articleID).Scan(&count)
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

	// Get uploaded file
	file, header, err := c.Request.FormFile("cover")
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "No file uploaded",
		})
		return
	}
	defer file.Close()

	// Upload to S3
	result, err := h.s3Service.UploadArticleCover(file, header, articleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to upload file: " + err.Error(),
		})
		return
	}

	// Update article cover image in database
	_, err = h.db.PostgreSQL.Exec(
		"UPDATE articles SET cover_image = $1 WHERE id = $2",
		result.URL, articleID,
	)
	if err != nil {
		// If database update fails, try to delete the uploaded file
		h.s3Service.DeleteFile(result.Key)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to update article cover",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Article cover uploaded successfully",
		"url":        result.URL,
		"article_id": articleID,
		"size":       result.Size,
	})
}

// UploadEBookFile uploads ebook file
func (h *UploadHandler) UploadEBookFile(c *gin.Context) {
	// Get uploaded file
	file, header, err := c.Request.FormFile("ebook")
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "No file uploaded",
		})
		return
	}
	defer file.Close()

	// Upload to S3
	result, err := h.s3Service.UploadEBook(file, header)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to upload file: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "EBook uploaded successfully",
		"url":       result.URL,
		"key":       result.Key,
		"size":      result.Size,
		"mime_type": result.MimeType,
	})
}

// UploadHarvestProof uploads harvest proof photo
func (h *UploadHandler) UploadHarvestProof(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: "User not authenticated",
		})
		return
	}

	// Get uploaded file
	file, header, err := c.Request.FormFile("proof")
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "No file uploaded",
		})
		return
	}
	defer file.Close()

	// TODO: Implement actual S3 upload when service is ready
	c.JSON(http.StatusOK, gin.H{
		"message":     "File upload functionality will be implemented with S3 service",
		"user_id":     userID,
		"filename":    header.Filename,
		"size":        header.Size,
		"placeholder": true,
	})
}