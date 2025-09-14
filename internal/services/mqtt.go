package services

import (
	"encoding/json"
	"fmt"
	"log"
	"swiflet-backend/internal/config"
	"swiflet-backend/internal/database"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTService struct {
	client mqtt.Client
	db     *database.DB
	config *config.Config
}

// SensorData represents incoming sensor data from MQTT
type SensorData struct {
	InstallCode string  `json:"install_code"`
	Suhu        float64 `json:"suhu"`
	Kelembaban  float64 `json:"kelembaban"`
	Timestamp   string  `json:"timestamp,omitempty"`
}

func NewMQTTService(cfg *config.Config, db *database.DB) (*MQTTService, error) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(cfg.MQTT.Broker)
	opts.SetClientID(cfg.MQTT.ClientID)
	
	if cfg.MQTT.Username != "" {
		opts.SetUsername(cfg.MQTT.Username)
		opts.SetPassword(cfg.MQTT.Password)
	}

	// Production settings
	opts.SetAutoReconnect(true)
	opts.SetCleanSession(true)
	opts.SetConnectTimeout(30 * time.Second)
	opts.SetKeepAlive(60 * time.Second)
	opts.SetPingTimeout(10 * time.Second)
	opts.SetWriteTimeout(10 * time.Second)
	opts.SetMaxReconnectInterval(10 * time.Minute)
	opts.SetConnectRetryInterval(5 * time.Second)
	opts.SetConnectRetry(true)

	service := &MQTTService{
		config: cfg,
		db:     db,
	}

	// Set connection lost handler
	opts.SetConnectionLostHandler(func(client mqtt.Client, err error) {
		log.Printf("MQTT connection lost: %v. Will attempt to reconnect...", err)
	})

	// Set on connect handler
	opts.SetOnConnectHandler(func(client mqtt.Client) {
		log.Println("MQTT connected successfully")
		service.subscribe()
	})

	// Set reconnect handler
	opts.SetReconnectingHandler(func(client mqtt.Client, opts *mqtt.ClientOptions) {
		log.Println("MQTT attempting to reconnect...")
	})

	client := mqtt.NewClient(opts)
	service.client = client

	return service, nil
}

// Connect to MQTT broker with retry logic
func (s *MQTTService) Connect() error {
	log.Println("Attempting to connect to MQTT broker...")
	
	if token := s.client.Connect(); token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to connect to MQTT broker: %w", token.Error())
	}
	
	log.Println("MQTT broker connection established")
	return nil
}

// ConnectWithRetry attempts to connect with exponential backoff
func (s *MQTTService) ConnectWithRetry(maxRetries int) error {
	var lastErr error
	backoff := 1 * time.Second
	
	for i := 0; i < maxRetries; i++ {
		if err := s.Connect(); err != nil {
			lastErr = err
			log.Printf("MQTT connection attempt %d/%d failed: %v", i+1, maxRetries, err)
			
			if i < maxRetries-1 {
				log.Printf("Retrying in %v...", backoff)
				time.Sleep(backoff)
				backoff *= 2
				if backoff > 30*time.Second {
					backoff = 30 * time.Second
				}
			}
			continue
		}
		return nil
	}
	
	return fmt.Errorf("failed to connect to MQTT after %d attempts: %w", maxRetries, lastErr)
}

// Disconnect from MQTT broker
func (s *MQTTService) Disconnect() {
	s.client.Disconnect(250)
}

// Subscribe to sensor data topics
func (s *MQTTService) subscribe() {
	topic := s.config.MQTT.TopicSensor
	token := s.client.Subscribe(topic, 1, s.handleSensorData)
	if token.Wait() && token.Error() != nil {
		log.Printf("Failed to subscribe to topic %s: %v", topic, token.Error())
	} else {
		log.Printf("Subscribed to topic: %s", topic)
	}
}

// Handle incoming sensor data
func (s *MQTTService) handleSensorData(client mqtt.Client, msg mqtt.Message) {
	log.Printf("Received sensor data from topic %s: %s", msg.Topic(), string(msg.Payload()))

	var data SensorData
	if err := json.Unmarshal(msg.Payload(), &data); err != nil {
		log.Printf("Failed to unmarshal sensor data: %v", err)
		return
	}

	// Parse timestamp or use current time
	timestamp := time.Now()
	if data.Timestamp != "" {
		if parsedTime, err := time.Parse(time.RFC3339, data.Timestamp); err == nil {
			timestamp = parsedTime
		}
	}

	// Validate install_code exists in iot_devices table
	var count int
	err := s.db.PostgreSQL.QueryRow("SELECT COUNT(*) FROM iot_devices WHERE install_code = $1", data.InstallCode).Scan(&count)
	if err != nil {
		log.Printf("Failed to validate install_code: %v", err)
		return
	}

	if count == 0 {
		log.Printf("Invalid install_code: %s", data.InstallCode)
		return
	}

	// Insert sensor data into TimescaleDB
	_, err = s.db.TimescaleDB.Exec(`
		INSERT INTO sensors (install_code, suhu, kelembaban, timestamp)
		VALUES ($1, $2, $3, $4)
	`, data.InstallCode, data.Suhu, data.Kelembaban, timestamp)

	if err != nil {
		log.Printf("Failed to insert sensor data: %v", err)
		return
	}

	log.Printf("Sensor data saved: %s - Suhu: %.2f, Kelembaban: %.2f", 
		data.InstallCode, data.Suhu, data.Kelembaban)
}

// PublishControlCommand publishes control commands to devices
func (s *MQTTService) PublishControlCommand(installCode string, command interface{}) error {
	topic := fmt.Sprintf("control/%s/command", installCode)
	
	payload, err := json.Marshal(command)
	if err != nil {
		return fmt.Errorf("failed to marshal command: %w", err)
	}

	token := s.client.Publish(topic, 1, false, payload)
	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to publish command: %w", token.Error())
	}

	log.Printf("Control command sent to %s: %s", installCode, string(payload))
	return nil
}