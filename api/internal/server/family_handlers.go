package server

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/zachczx/cubby/api/internal/response"
	"github.com/zachczx/cubby/api/internal/user"
)

func (s *Service) GetUsersFamiliesHandler(w http.ResponseWriter, r *http.Request) {
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

	response.WriteJSON(r.Context(), w, families)
}

func (s *Service) DeleteFamilyMemberHandler(w http.ResponseWriter, r *http.Request) {
	f := r.PathValue("familyID")
	familyID, err := uuid.Parse(f)
	if err != nil {
		response.WriteError(r.Context(), w, err)
		return
	}

	m := r.PathValue("memberID")
	memberID, err := uuid.Parse(m)
	if err != nil {
		response.WriteError(r.Context(), w, err)
		return
	}

	userID, err := s.GetUserIDFromContext(r.Context())
	if err != nil {
		response.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if err := user.DeleteMember(s.DB, familyID, userID, memberID); err != nil {
		response.WriteError(r.Context(), w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Service) LeaveFamilyHandler(w http.ResponseWriter, r *http.Request) {
	fid := r.PathValue("familyID")
	familyID, err := uuid.Parse(fid)
	if err != nil {
		response.WriteError(r.Context(), w, err)
		return
	}

	userID, err := s.GetUserIDFromContext(r.Context())
	if err != nil {
		response.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if err := user.LeaveFamily(s.DB, familyID, userID); err != nil {
		response.WriteError(r.Context(), w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
