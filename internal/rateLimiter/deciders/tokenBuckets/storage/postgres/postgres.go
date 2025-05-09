package postgres

import (
	"database/sql"
	"embed"
	_ "embed" // для go:embed
	"fmt"
	_ "github.com/lib/pq" // postgres driver
	"github.com/pressly/goose/v3"
	"go-load-balancer/internal/config"
)

type BucketSettingsStorage struct {
	db *sql.DB
}

//go:embed migrations/*.sql
var migrations embed.FS // Встраиваем все миграции из папки migrations/ в бинарник

// New создает подключение к БД и применяет миграции
func New(pgConfig config.BucketSettingsDatabase) (*BucketSettingsStorage, error) {
	const op = "rateLimiter.deciders.tokenBuckets.storage.postgres.New"

	// Подключаемся к базе данных
	db, err := sql.Open("postgres",
		fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
			pgConfig.Host, pgConfig.Port, pgConfig.User, pgConfig.Database, pgConfig.Password))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Применяем миграции с помощью Goose
	if err := applyMigrations(db); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &BucketSettingsStorage{db: db}, nil
}

// applyMigrations применяет все миграции, встроенные в бинарник
func applyMigrations(db *sql.DB) error {
	// Устанавливаем базу миграций из embed.FS
	goose.SetBaseFS(migrations)

	// Применяем миграции с помощью Goose
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set Goose dialect: %w", err)
	}

	// Применяем все миграции
	if err := goose.Up(db, "migrations"); err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	return nil
}
