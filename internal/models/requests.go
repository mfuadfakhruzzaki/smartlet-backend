package models

import (
	"time"
)

// InstallationRequest represents the InstallationRequest table
type InstallationRequest struct {
	ID              int       `json:"id" db:"id"`
	SwifletHouseID  int       `json:"id_swiflet_house" db:"id_swiflet_house" validate:"required"`
	Floors          string    `json:"floors" db:"floors" validate:"required"`
	SensorCount     int       `json:"sensor_count" db:"sensor_count" validate:"required"`
	AppointmentDate time.Time `json:"appointment_date" db:"appointment_date" validate:"required"`
	Status          int       `json:"status" db:"status"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

// MaintenanceRequest represents the MaintenanceRequest table
type MaintenanceRequest struct {
	ID              int       `json:"id" db:"id"`
	DeviceID        int       `json:"id_device" db:"id_device" validate:"required"`
	Reason          string    `json:"reason" db:"reason" validate:"required"`
	AppointmentDate time.Time `json:"appointment_date" db:"appointment_date" validate:"required"`
	Status          int       `json:"status" db:"status"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

// UninstallationRequest represents the UninstallationRequest table
type UninstallationRequest struct {
	ID              int       `json:"id" db:"id"`
	DeviceID        int       `json:"id_device" db:"id_device" validate:"required"`
	Reason          string    `json:"reason" db:"reason" validate:"required"`
	AppointmentDate time.Time `json:"appointment_date" db:"appointment_date" validate:"required"`
	Status          int       `json:"status" db:"status"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}