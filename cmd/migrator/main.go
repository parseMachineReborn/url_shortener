package main

import (
	"errors"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/joho/godotenv"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Проблема при загрузке .env")
	}

	connStr := os.Getenv("DB_CONNECTION_STRING")
	migrationsSrc := os.Getenv("MIGRATION_SOURCE_DIR")

	mg, err := migrate.New(migrationsSrc, connStr)
	if err != nil {
		log.Fatal("Ошибка при создании мигратора", err)
	}

	if err := mg.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return
		}
		log.Fatal("Ошибка при миграции", err)
	}
}
