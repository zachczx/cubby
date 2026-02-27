package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/zachczx/cubby/api/internal/notifier"
	"github.com/zachczx/cubby/api/internal/server"
	"github.com/zachczx/cubby/api/internal/tracker"
	"github.com/zachczx/cubby/api/internal/user"
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

	fmt.Println(os.Getenv("ENV"))

	ctx := context.Background()
	fcm, err := notifier.NewFCMClient(ctx)
	if err != nil {
		log.Fatal(err)
	}

	s := server.NewService(
		os.Getenv("STYTCH_PROJECT_ID"),
		os.Getenv("STYTCH_SECRET"),
		db,
		tracker.DefaultService{},
		user.UserManager{},
		fcm,
	)

	mux := NewHTTPHandler(s)

	fmt.Println("Listening on:", ":"+os.Getenv("API_LISTEN_ADDR"))
	server := &http.Server{
		Addr:              ":" + os.Getenv("API_LISTEN_ADDR"),
		ReadHeaderTimeout: 5 * time.Second,
		Handler:           CORS(mux),
	}

	go tracker.StartNotifications(s.DB, s.Notifier)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
