package server

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/zachczx/cubby/api/internal/response"
	"github.com/zachczx/cubby/api/internal/timer"
)

func (s *Service) GetAllTimerProfilesHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := s.GetUserIDFromContext(r.Context())
	if err != nil {
		response.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	profiles, err := timer.GetAllProfiles(s.DB, userID)
	if err != nil {
		response.WriteError(r.Context(), w, err)
		return
	}

	response.WriteJSON(r.Context(), w, profiles)
}

func (s *Service) NewTimerProfileHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := s.GetUserIDFromContext(r.Context())
	if err != nil {
		response.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var input timer.ProfileInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.WriteError(r.Context(), w, err)
		return
	}

	profile, err := timer.NewProfile(s.DB, userID, input)
	if err != nil {
		response.WriteError(r.Context(), w, err)
		return
	}

	response.WriteJSONStatus(r.Context(), w, http.StatusCreated, profile)
}

func (s *Service) EditTimerProfileHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := s.GetUserIDFromContext(r.Context())
	if err != nil {
		response.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	profileID, err := uuid.Parse(r.PathValue("profileID"))
	if err != nil {
		response.WriteError(r.Context(), w, err)
		return
	}

	var input timer.ProfileInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.WriteError(r.Context(), w, err)
		return
	}

	if err := timer.EditProfile(s.DB, userID, profileID, input); err != nil {
		response.WriteError(r.Context(), w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Service) DeleteTimerProfileHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := s.GetUserIDFromContext(r.Context())
	if err != nil {
		response.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	profileID, err := uuid.Parse(r.PathValue("profileID"))
	if err != nil {
		response.WriteError(r.Context(), w, err)
		return
	}

	if err := timer.DeleteProfile(s.DB, userID, profileID); err != nil {
		response.WriteError(r.Context(), w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
