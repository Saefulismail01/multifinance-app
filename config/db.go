package config

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func ConnectDB() (*sqlx.DB, *Config, error) {
	cfg, err := NewConfig()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load config: %w", err)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
	)

	db, err := sqlx.Connect(cfg.Driver, dsn)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to database: %w", err)
	}


	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)


	if err := db.Ping(); err != nil {
		return nil, nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, cfg, nil
}
