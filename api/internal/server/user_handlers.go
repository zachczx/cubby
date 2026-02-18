package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/zachczx/cubby/api/internal/response"
	"github.com/zachczx/cubby/api/internal/user"
)

func GetUsersFamiliesHandler(s *Service, db *sqlx.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error

		userID, err := s.GetUserIDFromContext(r.Context())
		if err != nil {
			response.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		families, err := user.GetUsersFamilies(db, userID)
		if err != nil {
			response.WriteError(w, err)
			return
		}

		response.WriteJSON(w, families)
	})
}

type TaskDays struct {
	TaskDays int `json:"taskDays"`
}

func (s *Service) ChangeTaskLookaheadDaysHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := s.GetUserIDFromContext(r.Context())
	if err != nil {
		response.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var days TaskDays

	if err := json.NewDecoder(r.Body).Decode(&days); err != nil {
		response.WriteError(w, err)
		return
	}

	if err := user.ChangeTaskLookaheadDays(s.DB, userID, days.TaskDays); err != nil {
		response.WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type SoundInput struct {
	SoundOn bool `json:"soundOn"`
}

func (s *Service) ToggleSoundHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := s.GetUserIDFromContext(r.Context())
	if err != nil {
		response.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var soundInput SoundInput

	if err := json.NewDecoder(r.Body).Decode(&soundInput); err != nil {
		response.WriteError(w, err)
		return
	}

	if err := user.ToggleSound(s.DB, userID, soundInput.SoundOn); err != nil {
		response.WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Service) CreateVacationHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := s.GetUserIDFromContext(r.Context())
	if err != nil {
		response.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	ownedFamilyID, err := user.GetUserFamilyID(s.DB, userID)
	if err != nil {
		response.WriteError(w, err)
		return
	}

	var input user.VacationRequest

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.WriteError(w, err)
		return
	}

	if err := validateVacationInputDateTimes(input); err != nil {
		response.WriteError(w, err)
		return
	}

	if err := user.CreateVacation(s.DB, userID, ownedFamilyID, input); err != nil {
		response.WriteError(w, err)
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
		response.WriteError(w, err)
		return
	}

	vacations, err := user.GetVacations(s.DB, families)
	if err != nil {
		response.WriteError(w, err)
		return
	}

	response.WriteJSON(w, vacations)
}

func (s *Service) DeleteVacationHandler(w http.ResponseWriter, r *http.Request) {
	vid := r.PathValue("vacationID")
	vacationID, err := uuid.Parse(vid)
	if err != nil {
		response.WriteError(w, err)
		return
	}

	userID, err := s.GetUserIDFromContext(r.Context())
	if err != nil {
		response.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if err := user.DeleteVacation(s.DB, userID, vacationID); err != nil {
		response.WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Service) GetUsersFamiliesHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := s.GetUserIDFromContext(r.Context())
	if err != nil {
		response.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	families, err := user.GetUsersFamilies(s.DB, userID)
	if err != nil {
		response.WriteError(w, err)
		return
	}

	response.WriteJSON(w, families)
}

func (s *Service) DeleteFamilyMemberHandler(w http.ResponseWriter, r *http.Request) {
	f := r.PathValue("familyID")
	familyID, err := uuid.Parse(f)
	if err != nil {
		response.WriteError(w, err)
		return
	}

	m := r.PathValue("memberID")
	memberID, err := uuid.Parse(m)
	if err != nil {
		response.WriteError(w, err)
		return
	}

	userID, err := s.GetUserIDFromContext(r.Context())
	if err != nil {
		response.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if err := user.DeleteMember(s.DB, familyID, userID, memberID); err != nil {
		response.WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Service) GetFamilyInvitesHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := s.GetUserIDFromContext(r.Context())
	if err != nil {
		response.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	invites, err := user.GetFamilyInvites(s.DB, userID)
	if err != nil {
		response.WriteError(w, err)
		return
	}

	response.WriteJSON(w, invites)
}

func (s *Service) GetFamilyInviteHandler(w http.ResponseWriter, r *http.Request) {
	iid := r.PathValue("inviteID")
	inviteID, err := uuid.Parse(iid)
	if err != nil {
		response.WriteError(w, err)
		return
	}

	userID, err := s.GetUserIDFromContext(r.Context())
	if err != nil {
		response.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	invite, err := user.GetFamilyInvite(s.DB, userID, inviteID)
	if err != nil {
		response.WriteError(w, err)
		return
	}

	response.WriteJSON(w, invite)
}

func (s *Service) AcceptFamilyInviteHandler(w http.ResponseWriter, r *http.Request) {
	iid := r.PathValue("inviteID")
	inviteID, err := uuid.Parse(iid)
	if err != nil {
		response.WriteError(w, err)
		return
	}

	userID, err := s.GetUserIDFromContext(r.Context())
	if err != nil {
		response.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if err := user.AcceptFamilyInvite(s.DB, userID, inviteID); err != nil {
		response.WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Service) DeclineFamilyInviteHandler(w http.ResponseWriter, r *http.Request) {
	iid := r.PathValue("inviteID")
	inviteID, err := uuid.Parse(iid)
	if err != nil {
		response.WriteError(w, err)
		return
	}

	userID, err := s.GetUserIDFromContext(r.Context())
	if err != nil {
		response.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if err := user.DeclineFamilyInvite(s.DB, userID, inviteID); err != nil {
		response.WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Service) CreateFamilyInviteHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := s.GetUserIDFromContext(r.Context())
	if err != nil {
		response.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	ownedFamilyID, err := user.GetUserFamilyID(s.DB, userID)
	if err != nil {
		response.WriteError(w, err)
		return
	}

	var invite user.InviteRequest

	if err := json.NewDecoder(r.Body).Decode(&invite); err != nil {
		response.WriteError(w, err)
		return
	}

	if err := user.CreateFamilyInvite(s.DB, ownedFamilyID, invite.InviteeEmail); err != nil {
		response.WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *Service) LeaveFamilyHandler(w http.ResponseWriter, r *http.Request) {
	fid := r.PathValue("familyID")
	familyID, err := uuid.Parse(fid)
	if err != nil {
		response.WriteError(w, err)
		return
	}

	userID, err := s.GetUserIDFromContext(r.Context())
	if err != nil {
		response.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if err := user.LeaveFamily(s.DB, familyID, userID); err != nil {
		response.WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
