package user

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Family struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	OwnerID   uuid.UUID `json:"ownerId" db:"owner_id"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

func NewFamily(db *sqlx.DB, f Family) (uuid.UUID, error) {
	var createdID uuid.UUID

	newFamily := Family{
		Name:    f.Name,
		OwnerID: f.OwnerID,
	}

	query := `INSERT INTO families (name, owner_id) VALUES (:name, :owner_id) RETURNING id`

	rows, err := db.NamedQuery(query, newFamily)
	if err != nil {
		return createdID, err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&createdID); err != nil {
			return createdID, err
		}
	}

	return createdID, nil
}

func GetUserFamilyID(db *sqlx.DB, userID uuid.UUID) (uuid.UUID, error) {
	var familyID uuid.UUID

	q := `SELECT id FROM families WHERE owner_id=$1`

	if err := db.QueryRow(q, userID).Scan(&familyID); err != nil {
		if err == sql.ErrNoRows {
			return familyID, fmt.Errorf("no family: %w", err)
		}

		return familyID, fmt.Errorf("fetch family err: %w", err)
	}

	return familyID, nil
}
