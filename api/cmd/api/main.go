package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/zachczx/cubby/api/internal/server"
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

	s := server.NewService(
		os.Getenv("STYTCH_PROJECT_ID"),
		os.Getenv("STYTCH_SECRET"),
		db,
	)

	mux := MakeHTTPHandlers(s)

	fmt.Println("Listening on:", ":"+os.Getenv("API_LISTEN_ADDR"))
	server := &http.Server{
		Addr:              ":" + os.Getenv("API_LISTEN_ADDR"),
		ReadHeaderTimeout: 5 * time.Second,
		Handler:           CORS(mux),
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
