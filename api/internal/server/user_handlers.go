package server

import (
	"encoding/json"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/zachczx/cubby/api/internal/response"
	"github.com/zachczx/cubby/api/internal/user"
)

func GetUsersFamiliesHandler(s *Service, db *sqlx.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error

		u := s.GetAuthenticatedUser(w, r)
		email := u.Emails[0].Email

		userID, err := s.UserManager.GetInternalUserID(db, email)
		if err != nil {
			response.WriteError(w, err)
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
	u := s.GetAuthenticatedUser(w, r)
	if u == nil {
		response.RespondWithError(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	email := u.Emails[0].Email

	userID, err := s.UserManager.GetInternalUserID(s.DB, email)
	if err != nil {
		response.WriteError(w, err)
		return
	}

	var days int

	if err := json.NewDecoder(r.Body).Decode(&days); err != nil {
		response.WriteError(w, err)
		return
	}

	if err := user.ChangeTaskLookaheadDays(s.DB, userID, days); err != nil {
		response.WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type SoundInput struct {
	SoundOn bool `json:"soundOn"`
}

func (s *Service) ToggleSoundHandler(w http.ResponseWriter, r *http.Request) {
	u := s.GetAuthenticatedUser(w, r)
	if u == nil {
		response.RespondWithError(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	email := u.Emails[0].Email

	userID, err := s.UserManager.GetInternalUserID(s.DB, email)
	if err != nil {
		response.WriteError(w, err)
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
