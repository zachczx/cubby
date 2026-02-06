package entry

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func Create(db *sqlx.DB, e Entry) error {
	q := `INSERT INTO logs 
			(tracker_id, interval, interval_unit, performed_by, performed_at, remark) 
			VALUES
			($1, $2, $3, $4, $5, $6)`

	if _, err := db.Exec(q, e.TrackerID, e.Interval, e.IntervalUnit, e.PerformedBy, e.PerformedAt, e.Remark); err != nil {
		return fmt.Errorf("create entry sql: %w", err)
	}

	return nil
}

func GetAll(db *sqlx.DB, userID uuid.UUID) ([]Entry, error) {
	var entries []Entry

	q := `SELECT * FROM logs 
			WHERE performed_by=$1 
			ORDER BY performed_at DESC`

	if err := db.Select(&entries, q, userID); err != nil {
		return entries, fmt.Errorf("entry query: %w", err)
	}

	return entries, nil
}
