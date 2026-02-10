package entry

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func Create(db *sqlx.DB, e Entry) (Entry, error) {
	q := `INSERT INTO entries (tracker_id, interval, interval_unit, performed_by, performed_at, remark) 
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id, tracker_id, interval, interval_unit, performed_by, performed_at, remark, created_at, updated_at`

	var new Entry
	err := db.QueryRow(q, e.TrackerID,
		e.Interval,
		e.IntervalUnit,
		e.PerformedBy,
		e.PerformedAt,
		e.Remark).
		Scan(
			&new.ID,
			&new.TrackerID,
			&new.Interval,
			&new.IntervalUnit,
			&new.PerformedBy,
			&new.PerformedAt,
			&new.Remark,
			&new.CreatedAt,
			&new.UpdatedAt,
		)
	if err != nil {
		return Entry{}, fmt.Errorf("create entry sql: %w", err)
	}

	return new, nil
}

func GetAll(db *sqlx.DB, userID uuid.UUID) ([]Entry, error) {
	var entries []Entry

	q := `SELECT * FROM entries 
			WHERE performed_by=$1 
			ORDER BY performed_at DESC`

	if err := db.Select(&entries, q, userID); err != nil {
		return entries, fmt.Errorf("entry query: %w", err)
	}

	return entries, nil
}
