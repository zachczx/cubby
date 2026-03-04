package server

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/zachczx/cubby/api/internal/response"
	"github.com/zachczx/cubby/api/internal/user"
)

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
