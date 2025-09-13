package handlers

import (
	"net/http"
	"strconv"
	"swiflet-backend/internal/database"
	"swiflet-backend/internal/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type IoTHandler struct {
	db       *database.DB
	validate *validator.Validate
}

func NewIoTHandler(db *database.DB) *IoTHandler {
	return &IoTHandler{
		db:       db,
		validate: validator.New(),
	}
}

// ListSwifletHouses returns paginated list of swiflet houses
func (h *IoTHandler) ListSwifletHouses(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage := 10
	offset := (page - 1) * perPage

	var total int
	err := h.db.PostgreSQL.QueryRow("SELECT COUNT(*) FROM swiflet_houses").Scan(&total)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Database error",
		})
		return
	}

	rows, err := h.db.PostgreSQL.Query(`
		SELECT id, id_user, name, location, created_at
		FROM swiflet_houses
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

	var houses []models.SwifletHouse
	for rows.Next() {
		var house models.SwifletHouse
		err := rows.Scan(&house.ID, &house.UserID, &house.Name, &house.Location, &house.CreatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error: "Database error",
			})
			return
		}
		houses = append(houses, house)
	}

	c.JSON(http.StatusOK, gin.H{"data": houses})
}

// CreateSwifletHouse creates a new swiflet house
func (h *IoTHandler) CreateSwifletHouse(c *gin.Context) {
	var house models.SwifletHouse
	if err := c.ShouldBindJSON(&house); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid request body",
		})
		return
	}

	if err := h.validate.Struct(house); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Validation failed",
		})
		return
	}

	_, err := h.db.PostgreSQL.Exec(`
		INSERT INTO swiflet_houses (id_user, name, location, created_at)
		VALUES ($1, $2, $3, $4)
	`, house.UserID, house.Name, house.Location, time.Now())

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to create swiflet house",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Swiflet house created successfully"})
}

// ListIoTDevices returns paginated list of IoT devices
func (h *IoTHandler) ListIoTDevices(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage := 10
	offset := (page - 1) * perPage

	var total int
	err := h.db.PostgreSQL.QueryRow("SELECT COUNT(*) FROM iot_devices").Scan(&total)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Database error",
		})
		return
	}

	rows, err := h.db.PostgreSQL.Query(`
		SELECT id, id_swiflet_house, floor, install_code, status, created_at, updated_at
		FROM iot_devices
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

	var devices []models.IoTDevice
	for rows.Next() {
		var device models.IoTDevice
		err := rows.Scan(&device.ID, &device.SwifletHouseID, &device.Floor, 
			&device.InstallCode, &device.Status, &device.CreatedAt, &device.UpdatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error: "Database error",
			})
			return
		}
		devices = append(devices, device)
	}

	c.JSON(http.StatusOK, gin.H{"data": devices})
}

// CreateIoTDevice creates a new IoT device
func (h *IoTHandler) CreateIoTDevice(c *gin.Context) {
	var device models.IoTDevice
	if err := c.ShouldBindJSON(&device); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid request body",
		})
		return
	}

	if err := h.validate.Struct(device); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Validation failed",
		})
		return
	}

	_, err := h.db.PostgreSQL.Exec(`
		INSERT INTO iot_devices (id_swiflet_house, floor, install_code, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, device.SwifletHouseID, device.Floor, device.InstallCode, device.Status, time.Now(), time.Now())

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to create IoT device",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "IoT device created successfully"})
}

// ListSensors returns paginated list of sensor data
func (h *IoTHandler) ListSensors(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage := 10
	offset := (page - 1) * perPage

	var total int
	err := h.db.TimescaleDB.QueryRow("SELECT COUNT(*) FROM sensors").Scan(&total)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Database error",
		})
		return
	}

	rows, err := h.db.TimescaleDB.Query(`
		SELECT id, install_code, suhu, kelembaban, timestamp
		FROM sensors
		ORDER BY timestamp DESC
		LIMIT $1 OFFSET $2
	`, perPage, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Database error",
		})
		return
	}
	defer rows.Close()

	var sensors []models.Sensor
	for rows.Next() {
		var sensor models.Sensor
		err := rows.Scan(&sensor.ID, &sensor.InstallCode, &sensor.Suhu, 
			&sensor.Kelembaban, &sensor.Timestamp)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error: "Database error",
			})
			return
		}
		sensors = append(sensors, sensor)
	}

	c.JSON(http.StatusOK, gin.H{"data": sensors})
}