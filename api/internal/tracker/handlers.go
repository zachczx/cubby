package tracker

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/zachczx/cubby/api/internal/response"
	"github.com/zachczx/cubby/api/internal/server"
	"github.com/zachczx/cubby/api/internal/user"
)

type Tracker struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	Owner        uuid.UUID  `json:"-" db:"owner_id"`
	Family       uuid.UUID  `json:"familyId" db:"family_id"`
	Name         string     `json:"name" db:"name"`
	Display      string     `json:"display" db:"display"`
	Interval     int        `json:"interval" db:"interval"`
	IntervalUnit string     `json:"intervalUnit" db:"interval_unit"`
	Category     string     `json:"category" db:"category"`
	Kind         string     `json:"kind" db:"kind"`
	ActionLabel  string     `json:"actionLabel" db:"action_label"`
	Pinned       bool       `json:"pinned" db:"pinned"`
	Show         bool       `json:"show" db:"show"`
	Icon         string     `json:"icon" db:"icon"`
	StartDate    *time.Time `json:"startDate,omitempty" db:"start_date"`
	Cost         *float64   `json:"cost,omitempty" db:"cost"`
	CreatedAt    time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at" json:"updated_at"`

	FamilyName string `json:"familyName" db:"family_name"`
	IsOwner    bool   `json:"isOwner" db:"-"`
}

type TrackerInput struct {
	Name         string   `json:"name"`
	Display      string   `json:"display"`
	Interval     int      `json:"interval"`
	IntervalUnit string   `json:"intervalUnit"`
	Category     string   `json:"category"`
	Kind         string   `json:"kind"`
	ActionLabel  string   `json:"actionLabel"`
	Pinned       bool     `json:"pinned"`
	Show         bool     `json:"show"`
	Icon         string   `json:"icon"`
	StartDate    string   `json:"startDate"`
	Cost         *float64 `json:"cost"`
}

func CreateHandler(s *server.Service, db *sqlx.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error

		userID, err := s.GetUserIDFromContext(r.Context())
		if err != nil {
			response.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		familyID, err := user.GetUserFamilyID(db, userID)
		if err != nil {
			response.WriteError(w, err)
			return
		}

		var input TrackerInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			response.WriteError(w, err)
			return
		}

		var startDate *time.Time
		if input.StartDate != "" {
			sd, err := time.Parse(time.RFC3339, input.StartDate)
			if err != nil {
				response.WriteError(w, err)
				return
			}
			startDate = &sd
		}

		t := Tracker{
			Owner:        userID,
			Family:       familyID,
			Name:         input.Name,
			Display:      input.Display,
			Interval:     input.Interval,
			IntervalUnit: input.IntervalUnit,
			Category:     input.Category,
			Kind:         input.Kind,
			ActionLabel:  input.ActionLabel,
			Icon:         input.Icon,
			Pinned:       input.Pinned,
			Show:         input.Show,
			Cost:         input.Cost,
			StartDate:    startDate,
		}

		trackerID, err := New(db, t)
		if err != nil {
			response.WriteError(w, err)
			return
		}

		response.WriteJSON(w, trackerID)
	})
}

func EditHandler(s *server.Service, db *sqlx.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error

		userID, err := s.GetUserIDFromContext(r.Context())
		if err != nil {
			response.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		tid := r.PathValue("trackerID")
		trackerID, err := uuid.Parse(tid)
		if err != nil {
			response.WriteError(w, err)
			return
		}

		var input TrackerInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			response.WriteError(w, err)
			return
		}

		var startDate *time.Time
		if input.StartDate != "" {
			sd, err := time.Parse(time.RFC3339, input.StartDate)
			if err != nil {
				response.WriteError(w, err)
				return
			}
			startDate = &sd
		}

		t := Tracker{
			ID:           trackerID,
			Owner:        userID,
			Name:         input.Name,
			Display:      input.Display,
			Interval:     input.Interval,
			IntervalUnit: input.IntervalUnit,
			Category:     input.Category,
			Kind:         input.Kind,
			ActionLabel:  input.ActionLabel,
			Icon:         input.Icon,
			Pinned:       input.Pinned,
			Show:         input.Show,
			Cost:         input.Cost,
			StartDate:    startDate,
		}

		if err := Edit(db, t); err != nil {
			response.WriteError(w, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})
}

func DeleteHandler(s *server.Service, db *sqlx.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := r.PathValue("trackerID")
		trackerID, err := uuid.Parse(t)
		if err != nil {
			response.WriteError(w, err)
			return
		}

		userID, err := s.GetUserIDFromContext(r.Context())
		if err != nil {
			response.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		if err := Delete(db, trackerID, userID); err != nil {
			response.WriteError(w, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})
}

func GetHandler(s *server.Service, db *sqlx.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := r.PathValue("trackerID")
		trackerID, err := uuid.Parse(t)
		if err != nil {
			response.WriteError(w, err)
			return
		}

		userID, err := s.GetUserIDFromContext(r.Context())
		if err != nil {
			response.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		tracker, err := Get(db, trackerID, userID)
		if err != nil {
			response.WriteError(w, err)
			return
		}

		response.WriteJSON(w, tracker)
	})
}

func GetAllHandler(s *server.Service, db *sqlx.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := s.GetUserIDFromContext(r.Context())
		if err != nil {
			response.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		tracker, err := GetAll(db, userID)
		if err != nil {
			response.WriteError(w, err)
			return
		}

		response.WriteJSON(w, tracker)
	})
}

type TrackerToggle struct {
	Pinned bool `json:"pinned" db:"pinned"`
	Show   bool `json:"show" db:"show"`
}

func TogglePinHandler(s *server.Service, db *sqlx.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := r.PathValue("trackerID")
		trackerID, err := uuid.Parse(t)
		if err != nil {
			response.WriteError(w, err)
			return
		}

		userID, err := s.GetUserIDFromContext(r.Context())
		if err != nil {
			response.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		var toggle TrackerToggle

		if err := json.NewDecoder(r.Body).Decode(&toggle); err != nil {
			response.WriteError(w, err)
			return
		}

		if err := TogglePin(db, userID, trackerID, toggle.Pinned); err != nil {
			response.WriteError(w, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})
}

func ToggleShowHandler(s *server.Service, db *sqlx.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := r.PathValue("trackerID")
		trackerID, err := uuid.Parse(t)
		if err != nil {
			response.WriteError(w, err)
			return
		}

		userID, err := s.GetUserIDFromContext(r.Context())
		if err != nil {
			response.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		var toggle TrackerToggle

		if err := json.NewDecoder(r.Body).Decode(&toggle); err != nil {
			response.WriteError(w, err)
			return
		}

		if err := ToggleShow(db, userID, trackerID, toggle.Show); err != nil {
			response.WriteError(w, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})
}
