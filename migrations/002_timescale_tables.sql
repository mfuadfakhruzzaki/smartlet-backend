-- TimescaleDB database schema  
-- Create database: CREATE DATABASE swiflet_timeseries;
-- Enable TimescaleDB extension: CREATE EXTENSION IF NOT EXISTS timescaledb;

-- Sensors table (TimescaleDB hypertable)
CREATE TABLE IF NOT EXISTS sensors (
    id SERIAL,
    install_code VARCHAR(255) NOT NULL,
    suhu DECIMAL(5,2) NOT NULL,
    kelembaban DECIMAL(5,2) NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id, timestamp)
);

-- Convert to hypertable (time-series table)
SELECT create_hypertable('sensors', 'timestamp', if_not_exists => TRUE);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_sensors_install_code ON sensors(install_code);
CREATE INDEX IF NOT EXISTS idx_sensors_timestamp ON sensors(timestamp DESC);
CREATE INDEX IF NOT EXISTS idx_sensors_install_code_timestamp ON sensors(install_code, timestamp DESC);

-- Create retention policy (optional: keep data for 1 year)
-- SELECT add_retention_policy('sensors', INTERVAL '1 year');