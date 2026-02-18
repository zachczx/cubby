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

	FamilyName string `json:"familyName" db:"family_name"`
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

	var hasPendingInvite bool

	cQ := `SELECT EXISTS(SELECT 1 FROM invites WHERE invitee_id= $1 AND family_id = $2 AND status = $3)`

	if err := db.QueryRow(cQ, inviteeID, familyID, StatusPending).Scan(&hasPendingInvite); err != nil {
		return fmt.Errorf("invite check already invited: %w", err)
	}

	if hasPendingInvite {
		return fmt.Errorf("invite already sent to this user")
	}

	iQ := `INSERT INTO invites (family_id, invitee_id)
			VALUES ($1, $2)`

	if _, err := db.Exec(iQ, familyID, inviteeID); err != nil {
		return fmt.Errorf("create invite: %w", err)
	}

	return nil
}

func GetFamilyInvites(db *sqlx.DB, userID uuid.UUID) ([]Invite, error) {
	var invites []Invite

	q := `SELECT i.id, i.family_id, i.invitee_id, i.status, i.created_at, i.updated_at, families.name AS family_name 
			FROM invites i
			LEFT JOIN families ON i.family_id = families.id
			WHERE invitee_id = $1 AND status = $2`

	if err := db.Select(&invites, q, userID, StatusPending); err != nil {
		return invites, fmt.Errorf("get invite: %w", err)
	}

	return invites, nil
}

func GetFamilyInvite(db *sqlx.DB, userID uuid.UUID, inviteID uuid.UUID) (Invite, error) {
	var invite Invite

	q := `SELECT i.id, i.family_id, i.invitee_id, i.status, i.created_at, i.updated_at, families.name AS family_name 
			FROM invites i
			LEFT JOIN families ON i.family_id = families.id
			WHERE i.id = $1 AND i.invitee_id = $2`

	if err := db.Get(&invite, q, inviteID, userID); err != nil {
		return invite, fmt.Errorf("get invite: %w", err)
	}

	return invite, nil
}

func AcceptFamilyInvite(db *sqlx.DB, userID uuid.UUID, inviteID uuid.UUID) error {
	// Get family ID
	currentInvite, err := GetFamilyInvite(db, userID, inviteID)
	if err != nil {
		return fmt.Errorf("get family invite: %w", err)
	}

	tx, err := db.Beginx()
	if err != nil {
		return fmt.Errorf("accept family begin tx: %w", err)
	}
	defer tx.Rollback()

	insertQ := `INSERT INTO families_users (family_id, user_id) VALUES ($1, $2)`

	if _, err := db.Exec(insertQ, currentInvite.FamilyID, userID); err != nil {
		return fmt.Errorf("insert families_users: %w", err)
	}

	// Modify invite to mark completed
	updateQ := `UPDATE invites SET status = $1 WHERE invitee_id = $2 AND id = $3`

	if _, err := db.Exec(updateQ, StatusAccepted, userID, inviteID); err != nil {
		return fmt.Errorf("update invites: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("accept family commit tx: %w", err)
	}

	return nil
}

func DeclineFamilyInvite(db *sqlx.DB, userID uuid.UUID, inviteID uuid.UUID) error {
	q := `UPDATE invites SET status = $1 WHERE invitee_id = $2 AND id = $3`

	if _, err := db.Exec(q, StatusDeclined, userID, inviteID); err != nil {
		return fmt.Errorf("update invites: %w", err)
	}

	return nil
}
