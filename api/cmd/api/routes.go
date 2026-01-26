package main

import (
	"net/http"
	"os"

	"github.com/zachczx/cubby/api/internal/server"
)

func MakeHTTPHandlers(s *server.Service) http.Handler {
	mux := http.NewServeMux()

	// Main
	mux.HandleFunc("GET /{$}", Index)
	mux.HandleFunc("/magic-link", s.SendMagicLinkHandler)
	mux.HandleFunc("/authenticate", s.AuthenticateHandler)
	mux.HandleFunc("/check", s.CheckHandler)

	return mux
}

func Index(w http.ResponseWriter, _ *http.Request) {
	if _, err := w.Write([]byte("Cubby")); err != nil {
		server.WriteError(w, err)
		return
	}
}

func Healthcheck(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if os.Getenv("ENV") == "development" {
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		}
		if os.Getenv("ENV") != "development" {
			w.Header().Set("Access-Control-Allow-Origin", "https://cubby.dev")
		}
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
