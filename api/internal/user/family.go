package user

import (
	"database/sql"
	"fmt"
	"strings"
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

type FamilyResponse struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	IsOwner   bool      `json:"isOwner" db:"-"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
	Owner     User      `json:"owner" db:"owner"`
	Members   []User    `json:"members" db:"members"`
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

func GetUsersFamilies(db *sqlx.DB, userID uuid.UUID) ([]FamilyResponse, error) {
	var families []FamilyResponse

	q := `SELECT 
				f.id, 
				f.name, 
				f.owner_id, 
				f.created_at, 
				f.updated_at,
				u.name AS owner_name, 
				u.email AS owner_email
			FROM families AS f
			LEFT JOIN users AS u ON u.id = f.owner_id
			WHERE f.owner_id = $1 
			OR EXISTS (
				SELECT 1 
				FROM families_users fu_check 
				WHERE fu_check.family_id = f.id AND fu_check.user_id = $1
			)`

	fRows, err := db.Query(q, userID)
	if err != nil {
		return nil, fmt.Errorf("fetch families err: %w", err)
	}
	defer fRows.Close()

	var f FamilyResponse
	var owner User
	for fRows.Next() {
		err := fRows.Scan(&f.ID, &f.Name, &owner.ID, &f.CreatedAt, &f.UpdatedAt, &owner.Name, &owner.Email)
		if err != nil {
			return nil, fmt.Errorf("family scan err: %w", err)
		}

		f.IsOwner = userID == owner.ID

		if owner.Name != nil {
			owner.Name = ToPtr(cleanEmail(*owner.Name))
		}

		f.Owner = owner

		families = append(families, f)
	}

	if len(families) == 0 {
		return families, nil
	}

	// Get family members

	familyIDs := make([]uuid.UUID, 0, len(families))
	for _, f := range families {
		familyIDs = append(familyIDs, f.ID)
	}

	query, args, err := sqlx.In(`SELECT fu.family_id, users.id, users.email, users.name 
									FROM families_users fu
									LEFT JOIN users ON fu.user_id = users.id
									WHERE fu.family_id IN (?)`, familyIDs)

	query = db.Rebind(query)
	mRows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("fetch members: %w", err)
	}
	defer mRows.Close()

	membersByFamily := make(map[uuid.UUID][]User)

	for mRows.Next() {
		var familyID uuid.UUID
		var user User

		if err := mRows.Scan(&familyID, &user.ID, &user.Email, &user.Name); err != nil {
			return nil, fmt.Errorf("scan member: %w", err)
		}

		membersByFamily[familyID] = append(membersByFamily[familyID], user)
	}

	for i := range families {
		families[i].Members = membersByFamily[families[i].ID]

		if families[i].Members == nil {
			families[i].Members = []User{}
		}
	}

	return families, nil
}

func cleanEmail(email string) string {
	s := strings.Split(email, "@")
	var clean string

	if len(s) > 0 {
		emailName := s[0]

		maxLength := 11

		if len(emailName) <= maxLength {
			return emailName
		}

		clean = emailName[:maxLength]
	}

	return clean
}

func ToPtr[T any](v T) *T {
	return &v
}

func DeleteMember(db *sqlx.DB, familyID uuid.UUID, ownerID uuid.UUID, memberID uuid.UUID) error {
	q := `DELETE FROM families_users
			WHERE family_id IN (SELECT id FROM families WHERE owner_id = $1)
			AND family_id = $2
			AND user_id = $3`

	if _, err := db.Exec(q, ownerID, familyID, memberID); err != nil {
		return fmt.Errorf("delete member err: %w", err)
	}

	return nil
}
