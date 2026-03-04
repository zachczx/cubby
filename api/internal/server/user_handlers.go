package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/zachczx/cubby/api/internal/response"
	"github.com/zachczx/cubby/api/internal/user"
)

func (s *Service) GetUserIDFromContext(ctx context.Context) (uuid.UUID, error) {
	id, ok := ctx.Value(UserIDKey).(uuid.UUID)
	if !ok {
		return uuid.Nil, fmt.Errorf("user id not found in context")
	}
	return id, nil
}

func (s *Service) GetUserEmailFromContext(ctx context.Context) (string, error) {
	email, ok := ctx.Value(EmailKey).(string)
	if !ok {
		return "", fmt.Errorf("user email not found in context")
	}
	return email, nil
}

func (s *Service) CheckHandler(w http.ResponseWriter, r *http.Request) {
	u := s.GetAuthenticatedUser(w, r)

	if u == nil || u.UserID == "" {
		response.RespondWithError(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	fmt.Println("check")

	w.WriteHeader(http.StatusNoContent)
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

type AccountInfoInput struct {
	Name       string `json:"name"`
	FamilyName string `json:"familyName"`
}

func (s *Service) UpdateAccountInfoHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := s.GetUserIDFromContext(r.Context())
	if err != nil {
		response.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var input AccountInfoInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.WriteError(w, err)
		return
	}

	if input.Name != "" {
		if err := user.UpdateName(s.DB, userID, input.Name); err != nil {
			response.WriteError(w, err)
			return
		}
	}

	if input.FamilyName != "" {
		ownedFamilyID, err := user.GetUserFamilyID(s.DB, userID)
		if err != nil {
			response.WriteError(w, err)
			return
		}

		if err := user.UpdateFamilyName(s.DB, ownedFamilyID, input.FamilyName); err != nil {
			response.WriteError(w, err)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
