package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stytchauth/stytch-go/v16/stytch/consumer/stytchapi"
	"github.com/zachczx/cubby/api/internal/notifier"
	"github.com/zachczx/cubby/api/internal/user"
)

type Service struct {
	Client                *stytchapi.API
	DB                    *sqlx.DB
	TrackerDefaultCreator TrackerDefaultCreator
	UserManager           UserManager
	Notifier              *notifier.FCMClient
	CookieConfig          CookieConfig
}

type TrackerDefaultCreator interface {
	CreateDefaults(db *sqlx.DB, userID uuid.UUID) error
}

type UserManager interface {
	GetInternalUserID(db *sqlx.DB, email string) (uuid.UUID, error)
	SyncUserInternal(db *sqlx.DB, email string, createdAt time.Time) (bool, uuid.UUID, error)
	Get(db *sqlx.DB, email string) (user.User, error)
}

func NewService(projectID string, secret string, DB *sqlx.DB, dc TrackerDefaultCreator, um UserManager, fcm *notifier.FCMClient, cc CookieConfig) *Service {
	client, err := stytchapi.NewClient(projectID, secret)
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	return &Service{
		Client:                client,
		DB:                    DB,
		TrackerDefaultCreator: dc,
		UserManager:           um,
		Notifier:              fcm,
		CookieConfig:          cc,
	}
}

type CookieConfig struct {
	Secure      bool
	SameSite    http.SameSite
	Partitioned bool
	Domain      string
	Path        string
	HTTPOnly    bool
}

func (s *Service) setSessionCookies(w http.ResponseWriter, jwt string, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:        "stytch_session_jwt",
		Value:       jwt,
		Path:        s.CookieConfig.Path,
		Domain:      s.CookieConfig.Domain,
		HttpOnly:    s.CookieConfig.HTTPOnly,
		Secure:      s.CookieConfig.Secure,
		SameSite:    s.CookieConfig.SameSite,
		Partitioned: s.CookieConfig.Partitioned,
		MaxAge:      5 * 60, // 5 mins
	})

	http.SetCookie(w, &http.Cookie{
		Name:        "stytch_session_token",
		Value:       token,
		Path:        s.CookieConfig.Path,
		Domain:      s.CookieConfig.Domain,
		HttpOnly:    s.CookieConfig.HTTPOnly,
		Secure:      s.CookieConfig.Secure,
		SameSite:    s.CookieConfig.SameSite,
		Partitioned: s.CookieConfig.Partitioned,
		MaxAge:      24 * 30 * 60 * 60, // 30 days
	})
}

func NewCookieConfig() CookieConfig {
	c := CookieConfig{
		Path:        "/",
		HTTPOnly:    true,
		Domain:      os.Getenv("COOKIE_DOMAIN"),
		Secure:      true,
		SameSite:    http.SameSiteDefaultMode,
		Partitioned: true,
	}

	if os.Getenv("ENV") == "development" {
		c.SameSite = http.SameSiteNoneMode
	}

	return c
}

func (s *Service) getUser(userID string) (user.User, error) {
	var u user.User

	q := `SELECT id, email, name FROM users WHERE id=$1`

	if err := s.DB.QueryRow(q, userID).Scan(&u.ID, &u.Email, &u.Name); err != nil {
		if err == sql.ErrNoRows {
			return u, fmt.Errorf("%w", err)
		}
		return u, fmt.Errorf("error fetch user: %w", err)
	}

	return u, nil
}
