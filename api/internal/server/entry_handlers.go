package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/zachczx/cubby/api/internal/entry"
	"github.com/zachczx/cubby/api/internal/response"
)

func (s *Service) CreateEntryHandler(w http.ResponseWriter, r *http.Request) {
	trackerID, err := uuid.Parse(r.PathValue("trackerID"))
	if err != nil {
		response.WriteError(w, err)
		return
	}

	var input entry.EntryInput

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
		response.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	e := entry.Entry{
		TrackerID:    trackerID,
		PerformedBy:  userID,
		PerformedAt:  performedAt,
		Interval:     input.Interval,
		IntervalUnit: input.IntervalUnit,
		Remark:       input.Remark,
	}

	new, err := entry.Create(s.DB, e)
	if err != nil {
		response.WriteError(w, err)
		return
	}

	response.WriteJSONStatus(w, http.StatusCreated, new)
}

func (s *Service) GetAllEntriesHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := s.GetUserIDFromContext(r.Context())
	if err != nil {
		response.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	entries, err := entry.GetAll(s.DB, userID)
	if err != nil {
		response.WriteError(w, err)
		return
	}

	response.WriteJSON(w, entries)
}

func (s *Service) DeleteEntryHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := s.GetUserIDFromContext(r.Context())
	if err != nil {
		response.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	e := r.PathValue("entryID")
	entryID, err := uuid.Parse(e)
	if err != nil {
		response.WriteError(w, err)
		return
	}

	if err := entry.Delete(s.DB, userID, entryID); err != nil {
		response.WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Service) EditEntryHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := s.GetUserIDFromContext(r.Context())
	if err != nil {
		response.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	e := r.PathValue("entryID")
	entryID, err := uuid.Parse(e)
	if err != nil {
		response.WriteError(w, err)
		return
	}

	var input entry.EntryInput

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

	if err := entry.Edit(s.DB, userID, entryID, performedAt); err != nil {
		response.WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
