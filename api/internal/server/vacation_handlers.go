package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/zachczx/cubby/api/internal/response"
	"github.com/zachczx/cubby/api/internal/user"
)

func (s *Service) CreateVacationHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := s.GetUserIDFromContext(r.Context())
	if err != nil {
		response.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	ownedFamilyID, err := user.GetUserFamilyID(s.DB, userID)
	if err != nil {
		response.WriteError(r.Context(), w, err)
		return
	}

	var input user.VacationRequest

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.WriteError(r.Context(), w, err)
		return
	}

	if err := validateVacationInputDateTimes(input); err != nil {
		response.WriteError(r.Context(), w, err)
		return
	}

	if err := user.CreateVacation(s.DB, userID, ownedFamilyID, input); err != nil {
		response.WriteError(r.Context(), w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func validateVacationInputDateTimes(input user.VacationRequest) error {
	if input.StartDateTime.After(input.EndDateTime) {
		return errors.New("end time must be after start time")
	}

	return nil
}

func (s *Service) GetVacationsHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := s.GetUserIDFromContext(r.Context())
	if err != nil {
		response.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	families, err := user.GetUsersFamilies(s.DB, userID)
	if err != nil {
		response.WriteError(r.Context(), w, err)
		return
	}

	vacations, err := user.GetVacations(s.DB, families)
	if err != nil {
		response.WriteError(r.Context(), w, err)
		return
	}

	response.WriteJSON(r.Context(), w, vacations)
}

func (s *Service) DeleteVacationHandler(w http.ResponseWriter, r *http.Request) {
	vid := r.PathValue("vacationID")
	vacationID, err := uuid.Parse(vid)
	if err != nil {
		response.WriteError(r.Context(), w, err)
		return
	}

	userID, err := s.GetUserIDFromContext(r.Context())
	if err != nil {
		response.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if err := user.DeleteVacation(s.DB, userID, vacationID); err != nil {
		response.WriteError(r.Context(), w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
