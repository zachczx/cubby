package user

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Invite struct {
	ID        uuid.UUID    `json:"id" db:"id"`
	FamilyID  uuid.UUID    `json:"familyId" db:"family_id"`
	InviteeID uuid.UUID    `json:"inviteeId" db:"invitee_id"`
	Status    InviteStatus `json:"status" db:"status"`
	CreatedAt time.Time    `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time    `json:"updatedAt" db:"updated_at"`
}

type InviteRequest struct {
	FamilyID     uuid.UUID `json:"familyId"`
	InviteeEmail string    `json:"inviteeEmail"`
}

// For accepting or declining an invite
type InviteActionRequest struct {
	InviteID uuid.UUID    `json:"inviteId"`
	Status   InviteStatus `json:"status"` // "accepted" or "declined"
}

type InviteStatus string

const (
	StatusPending  InviteStatus = "pending"
	StatusAccepted InviteStatus = "accepted"
	StatusDeclined InviteStatus = "declined"
)

// Then in your logic, you can easily validate:
func (s InviteStatus) IsValid() bool {
	switch s {
	case StatusPending, StatusAccepted, StatusDeclined:
		return true
	}
	return false
}

func CreateFamilyInvite(db *sqlx.DB, familyID uuid.UUID, inviteeEmail string) error {
	var inviteeID uuid.UUID

	sQ := `SELECT id FROM users WHERE email = $1`

	if err := db.QueryRow(sQ, inviteeEmail).Scan(&inviteeID); err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("invite no such user: %w", err)
		}

		return fmt.Errorf("invite get userID: %w", err)
	}

	iQ := `INSERT INTO invites (family_id, invitee_id)
			VALUES ($1, $2)`

	if _, err := db.Exec(iQ, familyID, inviteeID); err != nil {
		return fmt.Errorf("create invite: %w", err)
	}

	return nil
}
