package tracker

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func New(db *sqlx.DB, t Tracker) (uuid.UUID, error) {
	var newID uuid.UUID

	q := `INSERT INTO trackers (
				user_id, family_id, name, display, interval, interval_unit, 
				category, kind, action_label, pinned, show, icon, start_date, cost, 
				created_at, updated_at
			) VALUES (
				:user_id, :family_id, :name, :display, :interval, :interval_unit, 
				:category, :kind, :action_label, :pinned, :show, :icon, :start_date, :cost, 
				NOW(), NOW()
			) RETURNING id`

	rows, err := db.NamedQuery(q, t)
	if err != nil {
		return newID, err
	}
	defer rows.Close()

	if rows.Next() {
		rows.Scan(&newID)
	}

	return newID, nil
}

func Get(db *sqlx.DB, trackerID uuid.UUID, userID uuid.UUID) (Tracker, error) {
	var t Tracker
	fmt.Println(trackerID, userID)
	q := `SELECT * FROM trackers WHERE id=$1 AND user_id=$2`

	if err := db.Get(&t, q, trackerID, userID); err != nil {
		return Tracker{}, fmt.Errorf("select tracker: %w", err)
	}

	return t, nil
}

func GetAll(db *sqlx.DB, userID uuid.UUID) ([]Tracker, error) {
	var t []Tracker
	q := `SELECT t.*, f.name AS family_name FROM trackers t
			JOIN families f ON t.family_id = f.id   
			WHERE user_id=$1`

	if err := db.Select(&t, q, userID); err != nil {
		return nil, fmt.Errorf("select trackers: %w", err)
	}

	for i := range t {
		if userID == t[i].User {
			t[i].IsOwner = true
		}
	}

	return t, nil
}
