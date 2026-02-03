package db

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Store struct {
	gormDB *gorm.DB
	sqlDB  *sql.DB
}

type healthSQLDB interface {
	PingContext(ctx context.Context) error
	Close() error
	SetMaxOpenConns(n int)
	SetMaxIdleConns(n int)
	SetConnMaxLifetime(d time.Duration)
}

// SQLDB is a small abstraction that allows easier testing and controlled health checks.
type SQLDB interface {
	healthSQLDB
}

func New(databaseURL string, logger *slog.Logger) (*Store, error) {
	gormDB, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open gorm connection: %w", err)
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql db from gorm: %w", err)
	}

	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	logger.Info("database connected")
	return &Store{gormDB: gormDB, sqlDB: sqlDB}, nil
}

func (s *Store) Ping(ctx context.Context) error {
	if s == nil || s.sqlDB == nil {
		return fmt.Errorf("database not initialized")
	}
	return s.sqlDB.PingContext(ctx)
}

func (s *Store) Close() error {
	if s == nil || s.sqlDB == nil {
		return nil
	}
	return s.sqlDB.Close()
}

func (s *Store) Gorm() *gorm.DB {
	if s == nil {
		return nil
	}
	return s.gormDB
}
