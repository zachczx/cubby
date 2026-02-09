package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/stytchauth/stytch-go/v16/stytch/consumer/magiclinks"
	emailML "github.com/stytchauth/stytch-go/v16/stytch/consumer/magiclinks/email"
	"github.com/stytchauth/stytch-go/v16/stytch/consumer/sessions"
	"github.com/stytchauth/stytch-go/v16/stytch/consumer/users"
	"github.com/zachczx/cubby/api/internal/response"
	"github.com/zachczx/cubby/api/internal/user"

	"github.com/google/uuid"
)

type contextKey string

const (
	UserIDKey contextKey = "userID"
	EmailKey  contextKey = "email"
)

func (s *Service) RequireAuthentication(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u := s.GetAuthenticatedUser(w, r)

		if u == nil || u.UserID == "" {
			response.RespondWithError(w, http.StatusForbidden, "unauthorized access")
			return
		}

		email := ""
		if len(u.Emails) > 0 {
			email = u.Emails[0].Email
		} else {
			resp, err := s.Client.Users.Get(r.Context(), &users.GetParams{
				UserID: u.UserID,
			})
			if err != nil {
				log.Printf("Error fetching user details: %v", err)
				response.RespondWithError(w, http.StatusInternalServerError, "internal server error")
				return
			}
			if len(resp.Emails) > 0 {
				email = resp.Emails[0].Email
			}
		}

		if email == "" {
			response.RespondWithError(w, http.StatusBadRequest, "email required")
			return
		}

		localUserID, err := s.UserManager.GetInternalUserID(s.DB, email)
		if err != nil {
			log.Printf("Error getting internal user: %v", err)
			response.RespondWithError(w, http.StatusUnauthorized, "user not found")
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, localUserID)
		ctx = context.WithValue(ctx, EmailKey, email)

		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

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

func (s *Service) SendMagicLinkHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Printf("Error parsing form: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	email := r.Form.Get("email")
	if email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	_, err := s.Client.MagicLinks.Email.LoginOrCreate(
		r.Context(),
		&emailML.LoginOrCreateParams{
			Email: email,
		})
	if err != nil {
		log.Printf("Error sending email: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (s *Service) AuthenticateHandler(w http.ResponseWriter, r *http.Request) {
	tokenType := r.URL.Query().Get("stytch_token_type")
	token := r.URL.Query().Get("token")

	if tokenType != "magic_links" {
		fmt.Printf("Error: unrecognized token type %s\n", tokenType)
		http.Error(w, fmt.Sprintf("Unrecognized token type %s", tokenType), http.StatusBadRequest)
		return
	}

	resp, err := s.Client.MagicLinks.Authenticate(r.Context(), &magiclinks.AuthenticateParams{
		Token:                  token,
		SessionDurationMinutes: 43200,
	})
	if err != nil {
		fmt.Printf("Error authenticating: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set the session cookies (JWT + Refresh Token)
	s.setSessionCookies(w, resp.SessionJWT, resp.SessionToken)

	isNewUser, userID, err := s.UserManager.SyncUserInternal(s.DB, resp.User.Emails[0].Email, *resp.User.CreatedAt)
	if err != nil {
		fmt.Println(err)
		return
	}

	if isNewUser {
		if err := s.TrackerDefaultCreator.CreateDefaults(s.DB, userID); err != nil {
			fmt.Printf("Error creating default trackers: %v\n", err)
		}
	}

	redirectURL := os.Getenv("PUBLIC_WEB_URL")

	if isNewUser {
		redirectURL += "/profile/account?onboarding=true"
	}

	/*
		A HTML redirect is used here because http.Redirect() often causes page failure.
		Seems like a case of cookie not being set yet, but redirect succeding, therefore causing frontend to fail auth check.
	*/
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `<!DOCTYPE html>
				<html>
				<head>
					<meta http-equiv="refresh" content="0;url=%s">
					<script>window.location.href = "%s";</script>
				</head>
				<body>
					<p>Redirecting...</p>
				</body>
				</html>`,
		redirectURL,
		redirectURL)
}

func (s *Service) CheckHandler(w http.ResponseWriter, r *http.Request) {
	u := s.GetAuthenticatedUser(w, r)

	if u == nil || u.UserID == "" {
		response.RespondWithError(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Service) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	u := s.GetAuthenticatedUser(w, r)

	if u == nil || u.UserID == "" {
		response.RespondWithError(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	email := ""
	if len(u.Emails) > 0 {
		email = u.Emails[0].Email
	}

	if email == "" {
		response.RespondWithError(w, http.StatusUnauthorized, "email not found")
		return
	}

	localUser, err := s.UserManager.Get(s.DB, email)
	if err != nil {
		response.RespondWithError(w, http.StatusUnauthorized, "user not found")
		return
	}

	response.WriteJSON(w, localUser)
}

func (s *Service) GetAuthenticatedUser(w http.ResponseWriter, r *http.Request) *users.User {
	// JWT as first try
	if cookie, err := r.Cookie("stytch_session_jwt"); err == nil {
		if sess, err := s.Client.Sessions.AuthenticateJWTLocal(cookie.Value, 0); err == nil {
			return &users.User{
				UserID: sess.UserID,
			}
		}
	}

	// Session Token Fallback
	cookie, err := r.Cookie("stytch_session_token")
	if err != nil {
		return nil
	}

	resp, err := s.Client.Sessions.Authenticate(r.Context(), &sessions.AuthenticateParams{
		SessionToken:           cookie.Value,
		SessionDurationMinutes: 43200,
	})
	if err != nil {
		log.Printf("Error refreshing session: %v\n", err)
		return nil
	}

	s.setSessionCookies(w, resp.SessionJWT, resp.SessionToken)

	return &resp.User
}

func (s *Service) Logout(w http.ResponseWriter, r *http.Request) {
	expire := time.Now().Add(-7 * 24 * time.Hour)
	http.SetCookie(w, &http.Cookie{
		Name:    "stytch_session_jwt",
		Value:   "",
		Path:    "/",
		Expires: expire,
		MaxAge:  -1,
	})
	http.SetCookie(w, &http.Cookie{
		Name:    "stytch_session_token",
		Value:   "",
		Path:    "/",
		Expires: expire,
		MaxAge:  -1,
	})
}

func (s *Service) GetUsersFamiliesHandler(w http.ResponseWriter, r *http.Request) {
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

	if err := user.DeleteMember(s.DB, familyID, userID, memberID); err != nil {
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
