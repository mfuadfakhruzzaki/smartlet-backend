package services

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"swiflet-backend/internal/config"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
)

type S3Service struct {
	session  *session.Session
	uploader *s3manager.Uploader
	s3Client *s3.S3
	config   *config.Config
}

type UploadResult struct {
	URL      string `json:"url"`
	Key      string `json:"key"`
	Bucket   string `json:"bucket"`
	Size     int64  `json:"size"`
	MimeType string `json:"mime_type"`
}

// File type categories for validation
var (
	ImageTypes = []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
	DocTypes   = []string{".pdf", ".epub", ".mobi", ".doc", ".docx"}
	VideoTypes = []string{".mp4", ".avi", ".mov", ".wmv", ".flv"}
)

func NewS3Service(cfg *config.Config) (*S3Service, error) {
	// For MinIO, use default region if not specified
	region := cfg.S3.Region
	if region == "" {
		region = "us-east-1" // Default for MinIO
	}

	// Create AWS session (compatible with MinIO)
	awsConfig := &aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(cfg.S3.AccessKey, cfg.S3.SecretKey, ""),
	}
	
	// Configure for MinIO or other S3-compatible services
	if cfg.S3.Endpoint != "" {
		awsConfig.Endpoint = aws.String(cfg.S3.Endpoint)
		awsConfig.S3ForcePathStyle = aws.Bool(true) // Required for MinIO
		awsConfig.DisableSSL = aws.Bool(!strings.HasPrefix(cfg.S3.Endpoint, "https://")) // Auto-detect SSL
	}

	sess, err := session.NewSession(awsConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create S3 session: %w", err)
	}

	service := &S3Service{
		session:  sess,
		uploader: s3manager.NewUploader(sess),
		s3Client: s3.New(sess),
		config:   cfg,
	}

	// For MinIO, skip bucket creation check to avoid API port issues
	// Assume bucket exists and is accessible
	fmt.Printf("Using MinIO endpoint: %s\n", cfg.S3.Endpoint)
	fmt.Printf("Using bucket: %s\n", cfg.S3.Bucket)

	return service, nil
}

// UploadFile uploads a file to S3
func (s *S3Service) UploadFile(file multipart.File, header *multipart.FileHeader, folder string) (*UploadResult, error) {
	// Validate file type
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !s.isAllowedFileType(ext) {
		return nil, fmt.Errorf("file type %s is not allowed", ext)
	}

	// Validate file size (10MB limit for images, 50MB for docs)
	maxSize := int64(10 * 1024 * 1024) // 10MB default
	if s.isDocumentType(ext) {
		maxSize = int64(50 * 1024 * 1024) // 50MB for documents
	}

	if header.Size > maxSize {
		return nil, fmt.Errorf("file size %d exceeds maximum allowed size %d", header.Size, maxSize)
	}

	// Generate unique filename
	filename := s.generateUniqueFilename(header.Filename, folder)

	// Upload to S3
	result, err := s.uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(s.config.S3.Bucket),
		Key:         aws.String(filename),
		Body:        file,
		ContentType: aws.String(s.getMimeType(ext)),
		// Remove ACL setting as bucket doesn't support public ACLs
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	return &UploadResult{
		URL:      result.Location,
		Key:      filename,
		Bucket:   s.config.S3.Bucket,
		Size:     header.Size,
		MimeType: s.getMimeType(ext),
	}, nil
}

// UploadUserProfileImage uploads user profile image
func (s *S3Service) UploadUserProfileImage(file multipart.File, header *multipart.FileHeader, userID int) (*UploadResult, error) {
	folder := fmt.Sprintf("users/%d/profile", userID)
	return s.UploadFile(file, header, folder)
}

// UploadArticleCover uploads article cover image
func (s *S3Service) UploadArticleCover(file multipart.File, header *multipart.FileHeader, articleID int) (*UploadResult, error) {
	folder := fmt.Sprintf("articles/%d/cover", articleID)
	return s.UploadFile(file, header, folder)
}

// UploadEBook uploads ebook file
func (s *S3Service) UploadEBook(file multipart.File, header *multipart.FileHeader) (*UploadResult, error) {
	return s.UploadFile(file, header, "ebooks")
}

// UploadEBookThumbnail uploads ebook thumbnail
func (s *S3Service) UploadEBookThumbnail(file multipart.File, header *multipart.FileHeader, ebookID int) (*UploadResult, error) {
	folder := fmt.Sprintf("ebooks/%d/thumbnail", ebookID)
	return s.UploadFile(file, header, folder)
}

// UploadHarvestProof uploads harvest proof photo
func (s *S3Service) UploadHarvestProof(file multipart.File, header *multipart.FileHeader, userID int) (*UploadResult, error) {
	folder := fmt.Sprintf("harvests/%d/proof", userID)
	return s.UploadFile(file, header, folder)
}

// DeleteFile deletes a file from S3
func (s *S3Service) DeleteFile(key string) error {
	_, err := s.s3Client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(s.config.S3.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

// GeneratePresignedURL generates a presigned URL for secure file access
func (s *S3Service) GeneratePresignedURL(key string, expiration time.Duration) (string, error) {
	req, _ := s.s3Client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(s.config.S3.Bucket),
		Key:    aws.String(key),
	})

	url, err := req.Presign(expiration)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return url, nil
}

// Helper functions
func (s *S3Service) generateUniqueFilename(originalFilename, folder string) string {
	ext := filepath.Ext(originalFilename)
	name := strings.TrimSuffix(filepath.Base(originalFilename), ext)
	
	// Sanitize filename
	name = strings.ReplaceAll(name, " ", "_")
	name = strings.ToLower(name)
	
	// Add UUID for uniqueness
	uniqueID := uuid.New().String()[:8]
	timestamp := time.Now().Format("20060102_150405")
	
	if folder != "" {
		return fmt.Sprintf("%s/%s_%s_%s%s", folder, name, timestamp, uniqueID, ext)
	}
	return fmt.Sprintf("%s_%s_%s%s", name, timestamp, uniqueID, ext)
}

func (s *S3Service) isAllowedFileType(ext string) bool {
	allowedTypes := append(append(ImageTypes, DocTypes...), VideoTypes...)
	for _, allowedExt := range allowedTypes {
		if ext == allowedExt {
			return true
		}
	}
	return false
}

func (s *S3Service) isDocumentType(ext string) bool {
	for _, docExt := range DocTypes {
		if ext == docExt {
			return true
		}
	}
	return false
}

func (s *S3Service) getMimeType(ext string) string {
	mimeTypes := map[string]string{
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
		".webp": "image/webp",
		".pdf":  "application/pdf",
		".epub": "application/epub+zip",
		".mobi": "application/x-mobipocket-ebook",
		".doc":  "application/msword",
		".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		".mp4":  "video/mp4",
		".avi":  "video/x-msvideo",
		".mov":  "video/quicktime",
	}
	
	if mimeType, exists := mimeTypes[ext]; exists {
		return mimeType
	}
	return "application/octet-stream"
}

// ensureBucketExists creates the bucket if it doesn't exist (useful for MinIO)
func (s *S3Service) ensureBucketExists() error {
	// Check if bucket exists
	_, err := s.s3Client.HeadBucket(&s3.HeadBucketInput{
		Bucket: aws.String(s.config.S3.Bucket),
	})
	
	if err == nil {
		// Bucket exists
		return nil
	}
	
	// Bucket doesn't exist, create it
	_, err = s.s3Client.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(s.config.S3.Bucket),
	})
	
	if err != nil {
		return fmt.Errorf("failed to create bucket %s: %w", s.config.S3.Bucket, err)
	}
	
	return nil
}

// ListObjects lists all objects in the bucket
func (s *S3Service) ListObjects(prefix string) ([]*s3.Object, error) {
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(s.config.S3.Bucket),
	}
	
	if prefix != "" {
		input.Prefix = aws.String(prefix)
	}

	result, err := s.s3Client.ListObjectsV2(input)
	if err != nil {
		return nil, fmt.Errorf("failed to list objects: %w", err)
	}

	return result.Contents, nil
}

// SetBucketPolicyPublicRead sets bucket policy to allow public read access
func (s *S3Service) SetBucketPolicyPublicRead() error {
	bucketPolicy := fmt.Sprintf(`{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Principal": "*",
				"Action": ["s3:GetObject"],
				"Resource": ["arn:aws:s3:::%s/*"]
			}
		]
	}`, s.config.S3.Bucket)

	_, err := s.s3Client.PutBucketPolicy(&s3.PutBucketPolicyInput{
		Bucket: aws.String(s.config.S3.Bucket),
		Policy: aws.String(bucketPolicy),
	})

	if err != nil {
		return fmt.Errorf("failed to set bucket policy: %w", err)
	}

	return nil
}

// GetBucketPolicy gets current bucket policy
func (s *S3Service) GetBucketPolicy() (string, error) {
	result, err := s.s3Client.GetBucketPolicy(&s3.GetBucketPolicyInput{
		Bucket: aws.String(s.config.S3.Bucket),
	})

	if err != nil {
		return "", fmt.Errorf("failed to get bucket policy: %w", err)
	}

	if result.Policy == nil {
		return "", fmt.Errorf("no policy set")
	}

	return *result.Policy, nil
}