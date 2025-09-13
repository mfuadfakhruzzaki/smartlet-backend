package models

import (
	"time"
)

// WeeklyPrice represents the WeeklyPrice table
type WeeklyPrice struct {
	ID        int       `json:"id" db:"id"`
	Province  string    `json:"province" db:"province" validate:"required"`
	Price     float64   `json:"price" db:"price" validate:"required"`
	WeekStart time.Time `json:"week_start" db:"week_start" validate:"required"`
	WeekEnd   time.Time `json:"week_end" db:"week_end" validate:"required"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// HarvestSales represents the HarvestSales table
type HarvestSales struct {
	ID              int       `json:"id" db:"id"`
	UserID          int       `json:"id_user" db:"id_user" validate:"required"`
	Province        string    `json:"province" db:"province" validate:"required"`
	Price           float64   `json:"price" db:"price" validate:"required"`
	BowlWeight      float64   `json:"bowl_weight" db:"bowl_weight"`
	OvalWeight      float64   `json:"oval_weight" db:"oval_weight"`
	CornerWeight    float64   `json:"corner_weight" db:"corner_weight"`
	BrokenWeight    float64   `json:"broken_weight" db:"broken_weight"`
	AppointmentDate time.Time `json:"appointment_date" db:"appointment_date" validate:"required"`
	ProofPhoto      *string   `json:"proof_photo" db:"proof_photo"`
	Status          int       `json:"status" db:"status"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

// Transaction represents the Transaction table
type Transaction struct {
	ID              int       `json:"id" db:"id"`
	OrderID         string    `json:"order_id" db:"order_id" validate:"required"`
	Status          int       `json:"status" db:"status" validate:"required,oneof=0 1 2"`
	Amount          float64   `json:"amount" db:"amount" validate:"required"`
	PaymentType     string    `json:"payment_type" db:"payment_type" validate:"required"`
	TransactionTime time.Time `json:"transaction_time" db:"transaction_time"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

// Membership represents the Membership table
type Membership struct {
	ID       int       `json:"id" db:"id"`
	UserID   int       `json:"id_user" db:"id_user" validate:"required"`
	JoinDate time.Time `json:"join_date" db:"join_date" validate:"required"`
	ExpDate  time.Time `json:"exp_date" db:"exp_date" validate:"required"`
	OrderID  string    `json:"order_id" db:"order_id" validate:"required"`
	Status   int       `json:"status" db:"status"`
}