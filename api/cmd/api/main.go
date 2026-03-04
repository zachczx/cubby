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
	"github.com/zachczx/cubby/api/internal/database"
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

	db, err := sqlx.Connect("pgx", database.GetConnectionString())
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

	origins := []string{os.Getenv("CORS_DEV"), os.Getenv("CORS_WEB"), os.Getenv("CORS_DEV_ALT"), os.Getenv("CORS_PROD_APP")}

	s := server.NewService(
		os.Getenv("STYTCH_PROJECT_ID"),
		os.Getenv("STYTCH_SECRET"),
		db,
		tracker.DefaultService{},
		user.UserManager{},
		fcm,
		server.NewCookieConfig(),
		origins,
	)

	mux := NewHTTPHandler(s)

	fmt.Println("Listening on:", ":"+os.Getenv("API_LISTEN_ADDR"))
	server := &http.Server{
		Addr:              ":" + os.Getenv("API_LISTEN_ADDR"),
		ReadHeaderTimeout: 5 * time.Second,
		Handler:           CORS(s, mux),
	}

	go tracker.StartNotifications(s.DB, s.Notifier)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
