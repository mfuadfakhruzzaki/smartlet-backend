package models

import (
	"time"
)

// SwifletHouse represents the SwifletHouse table
type SwifletHouse struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"id_user" db:"id_user" validate:"required"`
	Name      string    `json:"name" db:"name" validate:"required"`
	Location  string    `json:"location" db:"location" validate:"required"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// IoTDevice represents the IoTDevice table
type IoTDevice struct {
	ID              int       `json:"id" db:"id"`
	SwifletHouseID  int       `json:"id_swiflet_house" db:"id_swiflet_house" validate:"required"`
	Floor           int       `json:"floor" db:"floor" validate:"required"`
	InstallCode     string    `json:"install_code" db:"install_code" validate:"required"`
	Status          int       `json:"status" db:"status"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

// Sensor represents the Sensor table (TimescaleDB)
type Sensor struct {
	ID          int       `json:"id" db:"id"`
	InstallCode string    `json:"install_code" db:"install_code" validate:"required"`
	Suhu        float64   `json:"suhu" db:"suhu" validate:"required"`
	Kelembaban  float64   `json:"kelembaban" db:"kelembaban" validate:"required"`
	Timestamp   time.Time `json:"timestamp" db:"timestamp"`
}

// Harvest represents the Harvest table
type Harvest struct {
	ID             int       `json:"id" db:"id"`
	UserID         int       `json:"id_user" db:"id_user" validate:"required"`
	SwifletHouseID int       `json:"id_swiflet_house" db:"id_swiflet_house" validate:"required"`
	Floor          int       `json:"floor" db:"floor" validate:"required"`
	BowlWeight     float64   `json:"bowl_weight" db:"bowl_weight"`
	BowlPieces     int       `json:"bowl_pieces" db:"bowl_pieces"`
	OvalWeight     float64   `json:"oval_weight" db:"oval_weight"`
	OvalPieces     int       `json:"oval_pieces" db:"oval_pieces"`
	CornerWeight   float64   `json:"corner_weight" db:"corner_weight"`
	CornerPieces   int       `json:"corner_pieces" db:"corner_pieces"`
	BrokenWeight   float64   `json:"broken_weight" db:"broken_weight"`
	BrokenPieces   int       `json:"broken_pieces" db:"broken_pieces"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}