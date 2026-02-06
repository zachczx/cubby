package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stytchauth/stytch-go/v16/stytch/consumer/stytchapi"
)

type Service struct {
	Client         *stytchapi.API
	DB             *sqlx.DB
	DefaultCreator DefaultCreator
}

type DefaultCreator interface {
	CreateDefaults(db *sqlx.DB, userID uuid.UUID) error
}

type User struct {
	ID            string     `db:"id" json:"id"`
	Email         string     `db:"email" json:"email"`
	PreferredName NullString `db:"preferred_name" json:"preferredName"`
	CreatedAt     time.Time  `db:"created_at" json:"createdAt"`
}

type NullString struct {
	sql.NullString
}

// This doesn't work with ns *NullString as a pointer receiver. Am reading the value, not modifying.
// Struct has a value, if it has a pointer, then *NullString would work.
func (ns NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}

	bytes, err := json.Marshal(ns.String)
	if err != nil {
		return []byte("null"), fmt.Errorf("marshal json: %w", err)
	}

	return bytes, nil
}

func (ns NullString) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		ns.Valid = false
		return nil
	}

	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("unmarshal json: %w", err)
	}
	ns.String = s
	ns.Valid = true
	return nil
}

func (u *User) NameString() string {
	if !u.PreferredName.Valid {
		return ""
	}
	return u.PreferredName.String
}

func NewService(projectID string, secret string, DB *sqlx.DB, dc DefaultCreator) *Service {
	client, err := stytchapi.NewClient(projectID, secret)
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	return &Service{
		Client:         client,
		DB:             DB,
		DefaultCreator: dc,
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

func (s *Service) getUser(userID string) (User, error) {
	var u User

	q := `SELECT id, email, preferred_name FROM users WHERE id=$1`

	if err := s.DB.QueryRow(q, userID).Scan(&u.ID, &u.Email, &u.PreferredName); err != nil {
		if err == sql.ErrNoRows {
			return u, fmt.Errorf("%w", err)
		}
		return u, fmt.Errorf("error fetch user: %w", err)
	}

	return u, nil
}
