package user

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Vacation struct {
	ID            string    `db:"id"            json:"id"`
	FamilyID      string    `db:"family_id"     json:"familyId"`
	CreatedBy     *string   `db:"created_by"    json:"createdBy"`
	StartDateTime time.Time `db:"start_date_time" json:"startDateTime"`
	EndDateTime   time.Time `db:"end_date_time"   json:"endDateTime"`
	Label         *string   `db:"label"         json:"label"`
	CreatedAt     time.Time `db:"created_at"    json:"createdAt"`
	UpdatedAt     time.Time `db:"updated_at"    json:"updatedAt"`
}

type VacationRequest struct {
	StartDateTime time.Time `json:"startDateTime"`
	EndDateTime   time.Time `json:"endDateTime"`
	Label         *string   `json:"label"`
}

func CreateVacation(db *sqlx.DB, userID uuid.UUID, familyID uuid.UUID, v VacationRequest) error {
	q := `INSERT INTO vacations (family_id, created_by, start_date_time, end_date_time, label) 
			VALUES ($1, $2, $3, $4, $5)`

	if _, err := db.Exec(q, familyID, userID, v.StartDateTime, v.EndDateTime, v.Label); err != nil {
		return fmt.Errorf("create vacation: %w", err)
	}

	return nil
}

func GetVacations(db *sqlx.DB, families []FamilyResponse) ([]Vacation, error) {
	var familyIDs []uuid.UUID

	for _, f := range families {
		familyIDs = append(familyIDs, f.ID)
	}

	var vacations []Vacation
	var err error

	query, args, err := sqlx.In("SELECT * FROM vacations WHERE family_id IN (?);", familyIDs)
	if err != nil {
		return nil, fmt.Errorf("vacations rebind: %w", err)
	}

	query = db.Rebind(query)
	if err = db.Select(&vacations, query, args...); err != nil {
		return nil, fmt.Errorf("select vacations: %w", err)
	}

	return vacations, nil
}

func DeleteVacation(db *sqlx.DB, userID uuid.UUID, vacationID uuid.UUID) error {
	q := `DELETE FROM vacations
			USING families 
			WHERE vacations.id = $1
			AND vacations.family_id = families.id
			AND families.owner_id = $2`

	if _, err := db.Exec(q, vacationID, userID); err != nil {
		return fmt.Errorf("delete entry: %w", err)
	}

	return nil
}
