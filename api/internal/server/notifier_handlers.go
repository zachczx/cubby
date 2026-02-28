package server

import (
	"encoding/json"
	"net/http"
	"slices"

	"github.com/zachczx/cubby/api/internal/notifier"
	"github.com/zachczx/cubby/api/internal/response"
)

type PushTokenInput struct {
	Token    string `json:"token"`
	Platform string `json:"platform"`
}

var platforms = []string{"web", "ios", "android"}

func (s *Service) PushTokenHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := s.GetUserIDFromContext(r.Context())
	if err != nil {
		response.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var t PushTokenInput

	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		response.WriteError(w, err)
		return
	}

	if !slices.Contains(platforms, t.Platform) {
		response.RespondWithError(w, http.StatusUnprocessableEntity, "invalid platform")
		return
	}

	if err := notifier.SavePushToken(s.DB, userID, t.Token, t.Platform); err != nil {
		response.WriteError(w, err)
		return
	}
}
