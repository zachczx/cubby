package user

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type User struct {
	ID        uuid.UUID `db:"id"         json:"id"`
	Email     string    `db:"email"      json:"email"`
	Name      *string   `db:"name"       json:"name"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

func SyncUserInternal(db *sqlx.DB, email string, createdAt time.Time) (bool, uuid.UUID, error) {
	var userID uuid.UUID
	newUser := false

	q := `SELECT id FROM users WHERE email=$1`

	if err := db.QueryRow(q, email).Scan(&userID); err != nil {
		if err != sql.ErrNoRows {
			return newUser, userID, fmt.Errorf("fetch user: %w", err)
		}

		qy := `INSERT INTO users (email, created_at) VALUES ($1, $2) RETURNING id`

		if err := db.QueryRow(qy, email, createdAt).Scan(&userID); err != nil {
			return newUser, userID, fmt.Errorf("insert user: %w", err)
		}

		newUser = true

		f := Family{Name: "Family", OwnerID: userID}
		if _, err := NewFamily(db, f); err != nil {
			return newUser, userID, fmt.Errorf("create family: %w", err)
		}
	}

	return newUser, userID, nil
}

func GetInternalUserID(db *sqlx.DB, email string) (uuid.UUID, error) {
	var userID uuid.UUID

	q := `SELECT id FROM users WHERE email=$1`

	if err := db.QueryRow(q, email).Scan(&userID); err != nil {
		if err == sql.ErrNoRows {
			return userID, fmt.Errorf("%w", err)
		}
		return userID, fmt.Errorf("error fetch user: %w", err)
	}

	return userID, nil
}

func Get(db *sqlx.DB, email string) (User, error) {
	var user User

	q := `SELECT * FROM users WHERE email=$1`

	if err := db.Get(&user, q, email); err != nil {
		return User{}, fmt.Errorf("error fetch user: %w", err)
	}

	return user, nil
}
