package entry

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Entry struct {
	ID           uuid.UUID `db:"id" json:"id"`
	TrackerID    uuid.UUID `db:"tracker_id" json:"trackerId"`
	Interval     int       `db:"interval" json:"interval"`
	IntervalUnit string    `db:"interval_unit" json:"intervalUnit"`
	PerformedBy  uuid.UUID `db:"performed_by" json:"performedBy"`
	PerformedAt  time.Time `db:"performed_at" json:"performedAt"`
	Remark       string    `db:"remark" json:"remark"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

type EntryInput struct {
	ID           *uuid.UUID `db:"id" json:"id"`
	TrackerID    uuid.UUID  `db:"tracker_id" json:"trackerId"`
	PerformedAt  *string    `db:"performed_at" json:"performedAt"`
	Interval     int        `db:"interval" json:"interval"`
	IntervalUnit string     `db:"interval_unit" json:"intervalUnit"`
	Remark       string     `db:"remark" json:"remark"`
}

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

func Delete(db *sqlx.DB, userID uuid.UUID, entryID uuid.UUID) error {
	q := `DELETE FROM entries
			USING trackers 
			WHERE entries.id = $1
			AND entries.tracker_id = trackers.id
			AND trackers.owner_id = $2`

	if _, err := db.Exec(q, entryID, userID); err != nil {
		return fmt.Errorf("delete entry: %w", err)
	}

	return nil
}

func Edit(db *sqlx.DB, userID uuid.UUID, entryID uuid.UUID, performedAt time.Time) error {
	q := `UPDATE entries
			SET performed_at = $1 
			FROM trackers
			WHERE trackers.id = entries.tracker_id
			AND trackers.owner_id = $2
			AND entries.id = $3`

	if _, err := db.Exec(q, performedAt, userID, entryID); err != nil {
		return fmt.Errorf("update entry: %w", err)
	}

	return nil
}
