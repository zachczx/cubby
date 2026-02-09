package server

import (
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
