package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"log"
	"usdt-rate-service/config"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env")
	}

	cfg := config.LoadConfig()

	var dbHost, dbPort, dbUser, dbPassword, dbName, migrationsPath string
	var down bool

	flag.StringVar(&dbHost, "db-host", cfg.DB.Host, "database host")
	flag.StringVar(&dbPort, "db-port", cfg.DB.Port, "database port")
	flag.StringVar(&dbUser, "db-user", cfg.DB.User, "database user")
	flag.StringVar(&dbPassword, "db-password", cfg.DB.Password, "database password")
	flag.StringVar(&dbName, "db-name", cfg.DB.Name, "database name")
	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations")
	flag.BoolVar(&down, "down", false, "Rollback Last migration")
	flag.Parse()

	if migrationsPath == "" {
		log.Fatal("migrations-path is required")
	}

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)

	m, err := migrate.New("file://"+migrationsPath, dbURL)
	if err != nil {
		log.Fatalf("failed to create migration instance: %v", err)
	}

	// Откат миграций, если передан флаг '-down'
	if down {
		if err := m.Down(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				fmt.Println("no migrations to rollback")
				return
			}
			log.Fatalf("failed to rollback migration: %v", err)
		}
		fmt.Println("Migrations rolled back successfully")
		return
	}

	// Выполняем миграции до последней версии
	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")
			return
		}
		log.Fatalf("failed to apply migrations: %v", err)
	}
	fmt.Println("migrations applied successfully")
}
