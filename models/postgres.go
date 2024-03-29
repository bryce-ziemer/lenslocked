package models

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"
)

// Open will open a SQL connection with the provided Postgres database.
// Callers of Open need to ensure that the connection is evenually closed via the db.Cklose() method
func Open(config PostgressConfig) (*sql.DB, error) {
	db, err := sql.Open("pgx", config.toString())
	if err != nil {
		return nil, fmt.Errorf("Open: %w", err)
	}

	return db, nil
}

func DefaultPostgresConfig() PostgressConfig {
	return PostgressConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "baloo",
		Password: "junglebook",
		Database: "lenslocked",
		SSLModel: "disable",
	}
}

type PostgressConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLModel string
}

func (cfg PostgressConfig) toString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.SSLModel)

}
