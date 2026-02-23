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
}

type TrackerDefaultCreator interface {
	CreateDefaults(db *sqlx.DB, userID uuid.UUID) error
}

type UserManager interface {
	GetInternalUserID(db *sqlx.DB, email string) (uuid.UUID, error)
	SyncUserInternal(db *sqlx.DB, email string, createdAt time.Time) (bool, uuid.UUID, error)
	Get(db *sqlx.DB, email string) (user.User, error)
}

func NewService(projectID string, secret string, DB *sqlx.DB, dc TrackerDefaultCreator, um UserManager, fcm *notifier.FCMClient) *Service {
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
	}
}

func (s *Service) setSessionCookies(w http.ResponseWriter, jwt string, token string) {
	secureOption := true
	sameSiteOption := http.SameSiteDefaultMode
	partitioned := true

	if os.Getenv("ENV") == "development" {
		secureOption = false
		sameSiteOption = http.SameSiteLaxMode
		partitioned = false
	}

	domain := os.Getenv("COOKIE_DOMAIN")

	http.SetCookie(w, &http.Cookie{
		Name:        "stytch_session_jwt",
		Value:       jwt,
		Path:        "/",
		Domain:      domain,
		HttpOnly:    true,
		Secure:      secureOption,
		SameSite:    sameSiteOption,
		Partitioned: partitioned,
		MaxAge:      5 * 60, // 5 mins
	})

	http.SetCookie(w, &http.Cookie{
		Name:        "stytch_session_token",
		Value:       token,
		Path:        "/",
		Domain:      domain,
		HttpOnly:    true,
		Secure:      secureOption,
		SameSite:    sameSiteOption,
		Partitioned: partitioned,
		MaxAge:      24 * 30 * 60 * 60, // 30 days
	})
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
