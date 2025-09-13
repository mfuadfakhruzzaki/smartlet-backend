package models

import (
	"time"
)

// User represents the Users table
type User struct {
	ID         int       `json:"id" db:"id"`
	Email      string    `json:"email" db:"email" validate:"required,email"`
	Name       string    `json:"name" db:"name" validate:"required"`
	Location   *string   `json:"location" db:"location"`
	NoTelp     *string   `json:"no_telp" db:"no_telp"`
	Password   string    `json:"password,omitempty" db:"password" validate:"required,min=6"`
	ImgProfile *string   `json:"img_profile" db:"img_profile"`
	Status     *int      `json:"status" db:"status"` // 0=inactive, 1=active, 2=suspended
	Role       *int      `json:"role" db:"role"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

// Article represents the Article table
type Article struct {
	ID         int       `json:"id" db:"id"`
	Title      string    `json:"title" db:"title" validate:"required"`
	Content    string    `json:"content" db:"content" validate:"required"`
	CoverImage *string   `json:"cover_image" db:"cover_image"`
	TagID      *int      `json:"tag_id" db:"tag_id"`
	Status     int       `json:"status" db:"status"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

// Tags represents the Tags table
type Tags struct {
	ID   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name" validate:"required"`
}

// Comment represents the Comment table
type Comment struct {
	ID        int       `json:"id" db:"id"`
	ArticleID int       `json:"article_id" db:"article_id" validate:"required"`
	UserID    int       `json:"user_id" db:"user_id" validate:"required"`
	Content   string    `json:"content" db:"content" validate:"required"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// EBook represents the EBook table
type EBook struct {
	ID            int       `json:"id" db:"id"`
	Title         string    `json:"title" db:"title" validate:"required"`
	FilePath      string    `json:"file_path" db:"file_path" validate:"required"`
	ThumbnailPath *string   `json:"thumbnail_path" db:"thumbnail_path"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

// Video represents the Video table
type Video struct {
	ID            int       `json:"id" db:"id"`
	Title         string    `json:"title" db:"title" validate:"required"`
	Description   *string   `json:"description" db:"description"`
	YoutubeLink   string    `json:"youtube_link" db:"youtube_link" validate:"required"`
	ThumbnailPath *string   `json:"thumbnail_path" db:"thumbnail_path"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}