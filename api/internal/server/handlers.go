package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/stytchauth/stytch-go/v16/stytch/consumer/magiclinks"
	emailML "github.com/stytchauth/stytch-go/v16/stytch/consumer/magiclinks/email"
	"github.com/stytchauth/stytch-go/v16/stytch/consumer/sessions"
	"github.com/stytchauth/stytch-go/v16/stytch/consumer/users"
)

func RespondWithError(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	fmt.Fprintln(w, message)
}

func (s *Service) RequireAuthentication(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u := s.getAuthenticatedUser(w, r)

		if u == nil || u.UserID == "" {
			RespondWithError(w, http.StatusForbidden, "unauthorized access")
			return
		}

		h.ServeHTTP(w, r)
	})
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

	s.syncUserInternal(resp.User.Emails[0].Email, *resp.User.CreatedAt)

	redirectURL := os.Getenv("PUBLIC_WEB_URL")

	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func (s *Service) CheckHandler(w http.ResponseWriter, r *http.Request) {
	user := s.getAuthenticatedUser(w, r)
	fmt.Println(user.Emails[0].Email)
}

func (s *Service) getAuthenticatedUser(w http.ResponseWriter, r *http.Request) *users.User {
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
