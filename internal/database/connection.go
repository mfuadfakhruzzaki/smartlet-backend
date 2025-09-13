package database

import (
	"database/sql"
	"fmt"
	"swiflet-backend/internal/config"

	_ "github.com/lib/pq"
)

type DB struct {
	PostgreSQL  *sql.DB
	TimescaleDB *sql.DB
}

// NewConnection creates new database connections
func NewConnection(cfg *config.Config) (*DB, error) {
	// PostgreSQL connection
	pgDSN := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	pgDB, err := sql.Open("postgres", pgDSN)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	if err := pgDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping PostgreSQL: %w", err)
	}

	// TimescaleDB connection
	tsDSN := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.TimescaleDB.Host,
		cfg.TimescaleDB.Port,
		cfg.TimescaleDB.User,
		cfg.TimescaleDB.Password,
		cfg.TimescaleDB.DBName,
		cfg.TimescaleDB.SSLMode,
	)

	tsDB, err := sql.Open("postgres", tsDSN)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to TimescaleDB: %w", err)
	}

	if err := tsDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping TimescaleDB: %w", err)
	}

	// Configure connection pools
	pgDB.SetMaxOpenConns(25)
	pgDB.SetMaxIdleConns(25)
	
	tsDB.SetMaxOpenConns(25)
	tsDB.SetMaxIdleConns(25)

	return &DB{
		PostgreSQL:  pgDB,
		TimescaleDB: tsDB,
	}, nil
}

// Close closes all database connections
func (db *DB) Close() error {
	if err := db.PostgreSQL.Close(); err != nil {
		return fmt.Errorf("failed to close PostgreSQL connection: %w", err)
	}
	
	if err := db.TimescaleDB.Close(); err != nil {
		return fmt.Errorf("failed to close TimescaleDB connection: %w", err)
	}
	
	return nil
}