package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/zachczx/cubby/api/internal/response"
	"github.com/zachczx/cubby/api/internal/tracker"
	"github.com/zachczx/cubby/api/internal/user"
)

func (s *Service) CreateHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	userID, err := s.GetUserIDFromContext(r.Context())
	if err != nil {
		response.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	familyID, err := user.GetUserFamilyID(s.DB, userID)
	if err != nil {
		response.WriteError(w, err)
		return
	}

	var input tracker.TrackerInput
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

	t := tracker.Tracker{
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

	trackerID, err := tracker.New(s.DB, t)
	if err != nil {
		response.WriteError(w, err)
		return
	}

	response.WriteJSON(w, trackerID)
}

func (s *Service) EditHandler(w http.ResponseWriter, r *http.Request) {
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

	var input tracker.TrackerInput
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

	t := tracker.Tracker{
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

	if err := tracker.Edit(s.DB, t); err != nil {
		response.WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Service) DeleteHandler(w http.ResponseWriter, r *http.Request) {
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

	if err := tracker.Delete(s.DB, trackerID, userID); err != nil {
		response.WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Service) GetHandler(w http.ResponseWriter, r *http.Request) {
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

	tracker, err := tracker.Get(s.DB, trackerID, userID)
	if err != nil {
		response.WriteError(w, err)
		return
	}

	response.WriteJSON(w, tracker)
}

func (s *Service) GetAllHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := s.GetUserIDFromContext(r.Context())
	if err != nil {
		response.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	tracker, err := tracker.GetAll(s.DB, userID)
	if err != nil {
		response.WriteError(w, err)
		return
	}

	response.WriteJSON(w, tracker)
}

type TrackerToggle struct {
	Pinned bool `json:"pinned" db:"pinned"`
	Show   bool `json:"show" db:"show"`
}

func (s *Service) TogglePinHandler(w http.ResponseWriter, r *http.Request) {
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

	if err := tracker.TogglePin(s.DB, userID, trackerID, toggle.Pinned); err != nil {
		response.WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Service) ToggleShowHandler(w http.ResponseWriter, r *http.Request) {
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

	if err := tracker.ToggleShow(s.DB, userID, trackerID, toggle.Show); err != nil {
		response.WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Service) GenerateHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := s.GetUserIDFromContext(r.Context())
	if err != nil {
		response.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	fmt.Println(userID)

	t, err := tracker.GetTrackersLast(s.DB)
	if err != nil {
		response.WriteError(w, err)
		return
	}

	newT, err := tracker.CalculateTrackersLastDue(s.DB, t)
	if err != nil {
		response.WriteError(w, err)
		return
	}

	response.WriteJSON(w, newT)
}
