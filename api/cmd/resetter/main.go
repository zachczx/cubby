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

	err = godotenv.Load(".env")
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}
	if err == nil {
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
