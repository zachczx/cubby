package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/zachczx/cubby/api/internal/migration"
)

func main() {
	var err error

	// Try loading .env file, but don't fail if it doesn't exist (e.g. in Docker)
	if err = godotenv.Load("../.env"); err != nil {
		if !os.IsNotExist(err) {
			log.Printf("Warning: Error loading .env file: %v", err)
		}
	} else {
		fmt.Println("env init: ok")
	}

	pg := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	db, err := sqlx.Connect("pgx", pg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	migration.WipeData(db)

	migration.Create(db)
}
