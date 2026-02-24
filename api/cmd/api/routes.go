package main

import (
	"net/http"
	"os"

	"github.com/zachczx/cubby/api/internal/response"
	"github.com/zachczx/cubby/api/internal/server"
)

func MakeHTTPHandlers(s *server.Service) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", Index)
	mux.HandleFunc("GET /health", Healthcheck)
	mux.HandleFunc("/magic-link", s.SendMagicLinkHandler)
	mux.HandleFunc("/authenticate", s.AuthenticateHandler)
	mux.Handle("GET /logout", http.HandlerFunc(s.Logout))

	mux.HandleFunc("GET /check", s.CheckHandler)
	mux.HandleFunc("GET /users", s.GetUserHandler)
	mux.HandleFunc("GET /users/me/families", s.RequireAuthentication(s.GetUsersFamiliesHandler))
	mux.HandleFunc("PATCH /users/me/account", s.RequireAuthentication(s.UpdateAccountInfoHandler))
	mux.HandleFunc("PATCH /users/me/sound", s.RequireAuthentication(s.ToggleSoundHandler))
	mux.HandleFunc("PATCH /users/me/task-lookahead", s.RequireAuthentication(s.ChangeTaskLookaheadDaysHandler))

	mux.HandleFunc("GET /families/invites", s.RequireAuthentication(s.GetFamilyInvitesHandler))
	mux.HandleFunc("GET /families/invites/{inviteID}", s.RequireAuthentication(s.GetFamilyInviteHandler))
	mux.HandleFunc("POST /families/invites", s.RequireAuthentication(s.CreateFamilyInviteHandler))
	mux.HandleFunc("POST /families/{familyID}/leave", s.RequireAuthentication(s.LeaveFamilyHandler))
	mux.HandleFunc("POST /families/invites/{inviteID}/accept", s.RequireAuthentication(s.AcceptFamilyInviteHandler))
	mux.HandleFunc("POST /families/invites/{inviteID}/decline", s.RequireAuthentication(s.DeclineFamilyInviteHandler))
	mux.HandleFunc("DELETE /families/{familyID}/{memberID}", s.RequireAuthentication(s.DeleteFamilyMemberHandler))

	mux.HandleFunc("GET /vacations", s.RequireAuthentication(s.GetVacationsHandler))
	mux.HandleFunc("POST /vacations", s.RequireAuthentication(s.CreateVacationHandler))
	mux.HandleFunc("DELETE /vacations/{vacationID}", s.RequireAuthentication(s.DeleteVacationHandler))

	mux.HandleFunc("GET /trackers", s.RequireAuthentication(s.GetAllHandler))
	mux.HandleFunc("GET /trackers/{trackerID}", s.RequireAuthentication(s.GetHandler))
	mux.HandleFunc("POST /trackers", s.RequireAuthentication(s.CreateHandler))
	mux.HandleFunc("POST /trackers/{trackerID}/entries", s.RequireAuthentication(s.CreateEntryHandler))
	mux.HandleFunc("PATCH /trackers/{trackerID}", s.RequireAuthentication(s.EditHandler))
	mux.HandleFunc("DELETE /trackers/{trackerID}", s.RequireAuthentication(s.DeleteHandler))
	mux.HandleFunc("PATCH /trackers/{trackerID}/pinned", s.RequireAuthentication(s.TogglePinHandler))
	mux.HandleFunc("PATCH /trackers/{trackerID}/show", s.RequireAuthentication(s.ToggleShowHandler))

	mux.HandleFunc("GET /entries", s.RequireAuthentication(s.GetAllEntriesHandler))
	mux.HandleFunc("DELETE /entries/{entryID}", s.RequireAuthentication(s.DeleteEntryHandler))
	mux.HandleFunc("PATCH /entries/{entryID}", s.RequireAuthentication(s.EditEntryHandler))

	mux.HandleFunc("GET /notifications", s.RequireAuthentication(s.NotificationHandler))
	mux.HandleFunc("POST /tokens", s.RequireAuthentication(s.PushTokenHandler))
	mux.HandleFunc("GET /notifications/generate", s.RequireAuthentication(s.GenerateHandler))

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
