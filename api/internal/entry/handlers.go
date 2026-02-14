package entry

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/zachczx/cubby/api/internal/response"
	"github.com/zachczx/cubby/api/internal/server"
)

type Entry struct {
	ID           uuid.UUID `db:"id" json:"id"`
	TrackerID    uuid.UUID `db:"tracker_id" json:"trackerId"`
	Interval     int       `db:"interval" json:"interval"`
	IntervalUnit string    `db:"interval_unit" json:"intervalUnit"`
	PerformedBy  uuid.UUID `db:"performed_by" json:"performedBy"`
	PerformedAt  time.Time `db:"performed_at" json:"performedAt"`
	Remark       string    `db:"remark" json:"remark"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

type EntryInput struct {
	ID           *uuid.UUID `db:"id" json:"id"`
	TrackerID    uuid.UUID  `db:"tracker_id" json:"trackerId"`
	PerformedAt  *string    `db:"performed_at" json:"performedAt"`
	Interval     int        `db:"interval" json:"interval"`
	IntervalUnit string     `db:"interval_unit" json:"intervalUnit"`
	Remark       string     `db:"remark" json:"remark"`
}

func CreateHandler(s *server.Service, db *sqlx.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		trackerID, err := uuid.Parse(r.PathValue("trackerID"))
		if err != nil {
			response.WriteError(w, err)
			return
		}

		var input EntryInput

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			response.WriteError(w, err)
			return
		}

		var performedAt time.Time
		if input.PerformedAt != nil {
			var err error

			performedAt, err = time.Parse(time.RFC3339, *input.PerformedAt)
			if err != nil {
				response.WriteError(w, err)
				return
			}
		} else {
			performedAt = time.Now()
		}

		userID, err := s.GetUserIDFromContext(r.Context())
		if err != nil {
			response.WriteError(w, err)
			return
		}

		e := Entry{
			TrackerID:    trackerID,
			PerformedBy:  userID,
			PerformedAt:  performedAt,
			Interval:     input.Interval,
			IntervalUnit: input.IntervalUnit,
			Remark:       input.Remark,
		}

		new, err := Create(db, e)
		if err != nil {
			response.WriteError(w, err)
			return
		}

		response.WriteJSONStatus(w, http.StatusCreated, new)
	})
}

func GetAllHandler(s *server.Service, db *sqlx.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := s.GetUserIDFromContext(r.Context())
		if err != nil {
			response.WriteError(w, err)
			return
		}

		entries, err := GetAll(db, userID)
		if err != nil {
			response.WriteError(w, err)
			return
		}

		response.WriteJSON(w, entries)
	})
}

func DeleteHandler(s *server.Service, db *sqlx.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := s.GetUserIDFromContext(r.Context())
		if err != nil {
			response.WriteError(w, err)
			return
		}

		e := r.PathValue("entryID")
		entryID, err := uuid.Parse(e)
		if err != nil {
			response.WriteError(w, err)
			return
		}

		if err := Delete(db, userID, entryID); err != nil {
			response.WriteError(w, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})
}

func EditHandler(s *server.Service, db *sqlx.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := s.GetUserIDFromContext(r.Context())
		if err != nil {
			response.WriteError(w, err)
			return
		}

		e := r.PathValue("entryID")
		entryID, err := uuid.Parse(e)
		if err != nil {
			response.WriteError(w, err)
			return
		}

		var input EntryInput

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			response.WriteError(w, err)
			return
		}

		var performedAt time.Time
		if input.PerformedAt != nil {
			var err error

			performedAt, err = time.Parse(time.RFC3339, *input.PerformedAt)
			if err != nil {
				response.WriteError(w, err)
				return
			}
		}

		if err := Edit(db, userID, entryID, performedAt); err != nil {
			response.WriteError(w, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})
}
