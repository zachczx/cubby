package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/zachczx/cubby/api/internal/database"
	"github.com/zachczx/cubby/api/internal/logging"
	"github.com/zachczx/cubby/api/internal/notifier"
	"github.com/zachczx/cubby/api/internal/server"
	"github.com/zachczx/cubby/api/internal/tracker"
	"github.com/zachczx/cubby/api/internal/user"
)

const (
	shutdownGrace  = 10 * time.Second
	defaultTimeout = 5 * time.Second
)

func main() {
	var err error

	// Try loading .env file, but don't fail if it doesn't exist (e.g. in Docker)
	if err = godotenv.Load("../.env"); err != nil {
		if !os.IsNotExist(err) {
			log.Printf("Warning: Error loading .env file: %v", err)
		}
	} else {
		slog.Info("App init", "env init", "ok")
	}

	logging.Init()

	db, err := sqlx.Connect("pgx", database.GetConnectionString())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	initCtx, cancel := context.WithTimeout(context.Background(), shutdownGrace)
	defer cancel()

	fcm, err := notifier.NewFCMClient(initCtx)
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

	slog.Info("server started", "port", ":"+os.Getenv("API_LISTEN_ADDR"))
	server := &http.Server{
		Addr:              ":" + os.Getenv("API_LISTEN_ADDR"),
		ReadHeaderTimeout: defaultTimeout,
		Handler:           RequestIDMiddleware(CORSMiddleware(s, mux)),
	}

	osCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := tracker.StartNotifications(osCtx, s.DB, s.Notifier); err != nil {
			slog.Error("notification failure", "error", err)
		}
	}()

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	<-osCtx.Done()
	slog.Info("server graceful shutdown started", "status", "started")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatal("server forced shutdown: ", err)
	}

	slog.Info("server exited", "status", "ok")
}
