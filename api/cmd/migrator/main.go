package main

import (
	"fmt"
	"log"
	"os"

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

	migration.Migrate()
}
