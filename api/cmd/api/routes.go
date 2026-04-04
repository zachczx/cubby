package main

import (
	"context"
	"net/http"
	"slices"

	"github.com/google/uuid"
	"github.com/zachczx/cubby/api/internal/logging"
	"github.com/zachczx/cubby/api/internal/response"
	"github.com/zachczx/cubby/api/internal/server"
)

func NewHTTPHandler(s *server.Service) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", Index)
	mux.HandleFunc("GET /health", Healthcheck)
	mux.HandleFunc("/magic-link", s.SendMagicLinkHandler)
	mux.HandleFunc("/authenticate", s.MagicLinkHandler)
	mux.HandleFunc("/otp/send", s.SendOTPHandler)
	mux.HandleFunc("/otp/verify", s.VerifyOTPHandler)
	mux.Handle("POST /logout", http.HandlerFunc(s.Logout))

	mux.HandleFunc("GET /check", s.CheckHandler)
	mux.HandleFunc("GET /users", s.GetUserHandler)
	mux.HandleFunc("GET /users/me/families", s.RequireAuthentication(s.GetUsersFamiliesHandler))
	mux.HandleFunc("PATCH /users/me/account", s.RequireAuthentication(s.UpdateAccountInfoHandler))
	mux.HandleFunc("PATCH /users/me/sound", s.RequireAuthentication(s.ToggleSoundHandler))
	mux.HandleFunc("PATCH /users/me/task-lookahead", s.RequireAuthentication(s.ChangeTaskLookaheadDaysHandler))
	mux.HandleFunc("PATCH /users/me/character", s.RequireAuthentication(s.ChangePreferredCharacterHandler))

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
	mux.HandleFunc("POST /trackers", s.RequireAuthentication(s.NewHandler))
	mux.HandleFunc("POST /trackers/{trackerID}/entries", s.RequireAuthentication(s.CreateEntryHandler))
	mux.HandleFunc("PATCH /trackers/{trackerID}", s.RequireAuthentication(s.EditHandler))
	mux.HandleFunc("DELETE /trackers/{trackerID}", s.RequireAuthentication(s.DeleteHandler))
	mux.HandleFunc("PATCH /trackers/{trackerID}/pinned", s.RequireAuthentication(s.TogglePinHandler))
	mux.HandleFunc("PATCH /trackers/{trackerID}/show", s.RequireAuthentication(s.ToggleShowHandler))
	mux.HandleFunc("POST /trackers/{trackerID}/toggle-mute", s.RequireAuthentication(s.ToggleMuteHandler))

	mux.HandleFunc("GET /entries", s.RequireAuthentication(s.GetAllEntriesHandler))
	mux.HandleFunc("DELETE /entries/{entryID}", s.RequireAuthentication(s.DeleteEntryHandler))
	mux.HandleFunc("PATCH /entries/{entryID}", s.RequireAuthentication(s.EditEntryHandler))

	mux.HandleFunc("GET /notifications/generate", s.RequireAuthentication(s.GenerateHandler))
	mux.HandleFunc("POST /tokens", s.RequireAuthentication(s.PushTokenHandler))

	mux.HandleFunc("GET /timer-profiles", s.RequireAuthentication(s.GetAllTimerProfilesHandler))
	mux.HandleFunc("POST /timer-profiles", s.RequireAuthentication(s.NewTimerProfileHandler))
	mux.HandleFunc("PATCH /timer-profiles/{profileID}", s.RequireAuthentication(s.EditTimerProfileHandler))
	mux.HandleFunc("DELETE /timer-profiles/{profileID}", s.RequireAuthentication(s.DeleteTimerProfileHandler))

	mux.HandleFunc("GET /gym/workouts", s.RequireAuthentication(s.GetAllWorkoutsHandler))
	mux.HandleFunc("POST /gym/workouts", s.RequireAuthentication(s.NewWorkoutHandler))
	mux.HandleFunc("PATCH /gym/workouts/{workoutID}", s.RequireAuthentication(s.EditWorkoutHandler))
	mux.HandleFunc("DELETE /gym/workouts/{workoutID}", s.RequireAuthentication(s.DeleteWorkoutHandler))
	mux.HandleFunc("POST /gym/workouts/{workoutID}/sets", s.RequireAuthentication(s.NewSetHandler))
	mux.HandleFunc("PATCH /gym/sets/{setID}", s.RequireAuthentication(s.EditSetHandler))
	mux.HandleFunc("POST /gym/sets/reorder", s.RequireAuthentication(s.ReorderSetHandler))
	mux.HandleFunc("DELETE /gym/sets/{setID}", s.RequireAuthentication(s.DeleteSetHandler))
	mux.HandleFunc("GET /gym/favourites", s.RequireAuthentication(s.GetFavouritesHandler))
	mux.HandleFunc("POST /gym/favourites", s.RequireAuthentication(s.ToggleFavouriteHandler))
	mux.HandleFunc("GET /gym/routines", s.RequireAuthentication(s.GetAllRoutinesHandler))
	mux.HandleFunc("POST /gym/routines", s.RequireAuthentication(s.NewRoutineHandler))
	mux.HandleFunc("PATCH /gym/routines/{routineID}", s.RequireAuthentication(s.EditRoutineHandler))
	mux.HandleFunc("POST /gym/routines/reorder", s.RequireAuthentication(s.ReorderRoutineHandler))
	mux.HandleFunc("DELETE /gym/routines/{routineID}", s.RequireAuthentication(s.DeleteRoutineHandler))
	mux.HandleFunc("POST /gym/routines/{routineID}/exercises", s.RequireAuthentication(s.AddRoutineExerciseHandler))
	mux.HandleFunc("PATCH /gym/routines/{routineID}/exercises/{exerciseID}", s.RequireAuthentication(s.EditRoutineExerciseHandler))
	mux.HandleFunc("DELETE /gym/routines/{routineID}/exercises/{exerciseID}", s.RequireAuthentication(s.RemoveRoutineExerciseHandler))
	mux.HandleFunc("POST /gym/routines/{routineID}/exercises/reorder", s.RequireAuthentication(s.ReorderRoutineExerciseHandler))
	mux.HandleFunc("POST /gym/routines/{routineID}/start", s.RequireAuthentication(s.StartWorkoutFromRoutineHandler))

	mux.HandleFunc("GET /gym/stats/summary", s.RequireAuthentication(s.GetGymSummaryHandler))
	mux.HandleFunc("GET /gym/stats/calendar", s.RequireAuthentication(s.GetGymCalendarHandler))
	mux.HandleFunc("GET /gym/stats/muscles", s.RequireAuthentication(s.GetGymMusclesHandler))
	mux.HandleFunc("GET /gym/stats/exercises", s.RequireAuthentication(s.GetGymUserExercisesHandler))
	mux.HandleFunc("GET /gym/stats/exercises/{exerciseID}", s.RequireAuthentication(s.GetGymExerciseStatsHandler))

	mux.HandleFunc("GET /market/prices", s.RequireAuthentication(s.GetMarketPricesHandler))
	mux.HandleFunc("GET /market/prices/{priceID}", s.RequireAuthentication(s.GetMarketPriceHandler))
	mux.HandleFunc("POST /market/prices", s.RequireAuthentication(s.LogMarketPriceHandler))
	mux.HandleFunc("PATCH /market/prices/{priceID}", s.RequireAuthentication(s.UpdateMarketPriceHandler))
	mux.HandleFunc("DELETE /market/prices/{priceID}", s.RequireAuthentication(s.DeleteMarketPriceHandler))
	mux.HandleFunc("GET /market/insights", s.RequireAuthentication(s.GetMarketInsightsHandler))

	return mux
}

func Index(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("Cubby")); err != nil {
		response.WriteError(r.Context(), w, err)
		return
	}
}

func Healthcheck(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func CORSMiddleware(s *server.Service, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		if slices.Contains(s.AllowedOrigins, origin) {
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

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, _ := uuid.NewV7()
		uuid := u.String()
		ctx := context.WithValue(r.Context(), logging.RequestIDKey, uuid)
		w.Header().Set("X-Request-ID", uuid)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
