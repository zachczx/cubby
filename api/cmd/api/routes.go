package main

import (
	"net/http"
	"os"

	"github.com/zachczx/cubby/api/internal/entry"
	"github.com/zachczx/cubby/api/internal/response"
	"github.com/zachczx/cubby/api/internal/server"
	"github.com/zachczx/cubby/api/internal/tracker"
)

func MakeHTTPHandlers(s *server.Service) http.Handler {
	mux := http.NewServeMux()

	// Main
	mux.HandleFunc("GET /{$}", Index)
	mux.HandleFunc("GET /health", Healthcheck)
	mux.HandleFunc("/magic-link", s.SendMagicLinkHandler)
	mux.HandleFunc("/authenticate", s.AuthenticateHandler)
	mux.Handle("GET /logout", http.HandlerFunc(s.Logout))

	mux.HandleFunc("/check", s.CheckHandler)
	mux.HandleFunc("GET /users", s.GetUserHandler)
	mux.Handle("GET /users/me/families", s.RequireAuthentication(http.HandlerFunc(s.GetUsersFamiliesHandler)))
	mux.Handle("DELETE /families/{familyID}/{memberID}", s.RequireAuthentication(http.HandlerFunc(s.DeleteFamilyMemberHandler)))
	mux.Handle("PATCH /users/me/sound", s.RequireAuthentication(http.HandlerFunc(s.ToggleSoundHandler)))
	mux.Handle("PATCH /users/me/task-lookahead", s.RequireAuthentication(http.HandlerFunc(s.ChangeTaskLookaheadDaysHandler)))

	mux.Handle("GET /vacations", s.RequireAuthentication(http.HandlerFunc(s.GetVacationsHandler)))
	mux.Handle("POST /vacations", s.RequireAuthentication(http.HandlerFunc(s.CreateVacationHandler)))
	mux.Handle("DELETE /vacations/{vacationID}", s.RequireAuthentication(http.HandlerFunc(s.DeleteVacationHandler)))

	mux.Handle("GET /trackers", s.RequireAuthentication(tracker.GetAllHandler(s, s.DB)))
	mux.Handle("GET /trackers/{trackerID}", s.RequireAuthentication(tracker.GetHandler(s, s.DB)))
	mux.Handle("POST /trackers", tracker.CreateHandler(s, s.DB))
	mux.Handle("POST /trackers/{trackerID}/entries", s.RequireAuthentication(entry.CreateHandler(s, s.DB)))
	mux.Handle("PATCH /trackers/{trackerID}", s.RequireAuthentication(tracker.EditHandler(s, s.DB)))
	mux.Handle("DELETE /trackers/{trackerID}", s.RequireAuthentication(tracker.DeleteHandler(s, s.DB)))

	mux.Handle("GET /entries", s.RequireAuthentication(entry.GetAllHandler(s, s.DB)))
	mux.Handle("DELETE /entries/{entryID}", s.RequireAuthentication(entry.DeleteHandler(s, s.DB)))
	mux.Handle("PATCH /entries/{entryID}", s.RequireAuthentication(entry.EditHandler(s, s.DB)))

	return mux
}

func Index(w http.ResponseWriter, _ *http.Request) {
	if _, err := w.Write([]byte("Cubby")); err != nil {
		response.WriteError(w, err)
		return
	}
}

func Healthcheck(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

var (
	devOrigin     = "http://localhost:5173"
	allowedOrigin = os.Getenv("PUBLIC_WEB_URL")
)

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		if origin == allowedOrigin || origin == devOrigin {
			w.Header().Set("Access-Control-Allow-Origin", origin)
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
