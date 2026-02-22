package tracker

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func New(db *sqlx.DB, t Tracker) (uuid.UUID, error) {
	var newID uuid.UUID

	q := `INSERT INTO trackers (
				owner_id, family_id, name, display, interval, interval_unit, 
				category, kind, action_label, pinned, show, icon, start_date, cost, 
				created_at, updated_at
			) VALUES (
				:owner_id, :family_id, :name, :display, :interval, :interval_unit, 
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

func Edit(db *sqlx.DB, t Tracker) error {
	q := `UPDATE trackers 
			SET name = :name, 
				display = :display, 
				interval = :interval, 
				interval_unit = :interval_unit, 
				category = :category, 
				kind = :kind, 
				action_label = :action_label, 
				pinned = :pinned, 
				show = :show, 
				icon = :icon, 
				start_date = :start_date, 
				cost = :cost, 
				updated_at = NOW()
			FROM families
			WHERE trackers.owner_id = families.owner_id
			AND trackers.id = :id`

	if _, err := db.NamedExec(q, t); err != nil {
		return fmt.Errorf("edit tracker: %w", err)
	}

	return nil
}

func Delete(db *sqlx.DB, trackerID uuid.UUID, userID uuid.UUID) error {
	q := `DELETE FROM trackers WHERE id = $1 AND owner_id = $2`

	if _, err := db.Exec(q, trackerID, userID); err != nil {
		return fmt.Errorf("delete tracker: %w", err)
	}

	return nil
}

func Get(db *sqlx.DB, trackerID uuid.UUID, userID uuid.UUID) (Tracker, error) {
	var t Tracker
	fmt.Println(trackerID, userID)
	q := `SELECT * FROM trackers WHERE id=$1 AND owner_id=$2`

	if err := db.Get(&t, q, trackerID, userID); err != nil {
		return Tracker{}, fmt.Errorf("select tracker: %w", err)
	}

	return t, nil
}

func GetAll(db *sqlx.DB, userID uuid.UUID) ([]Tracker, error) {
	var t []Tracker
	q := `SELECT t.*, f.name AS family_name FROM trackers t
			JOIN families f ON t.family_id = f.id   
			WHERE t.owner_id = $1 OR t.family_id IN (
				SELECT family_id FROM families_users WHERE user_id = $1
			)
			ORDER BY t.pinned DESC, t.name ASC`

	if err := db.Select(&t, q, userID); err != nil {
		return nil, fmt.Errorf("select trackers: %w", err)
	}

	for i := range t {
		if userID == t[i].Owner {
			t[i].IsOwner = true
		}
	}

	return t, nil
}

func TogglePin(db *sqlx.DB, userID uuid.UUID, trackerID uuid.UUID, isPinned bool) error {
	q := `UPDATE trackers 
			SET pinned = $1
			WHERE id = $2 AND owner_id = $3`

	if _, err := db.Exec(q, isPinned, trackerID, userID); err != nil {
		return fmt.Errorf("toggle pin: %w", err)
	}

	return nil
}

func ToggleShow(db *sqlx.DB, userID uuid.UUID, trackerID uuid.UUID, show bool) error {
	q := `UPDATE trackers 
			SET show = $1
			WHERE id = $2 AND owner_id = $3`

	if _, err := db.Exec(q, show, trackerID, userID); err != nil {
		return fmt.Errorf("toggle pin: %w", err)
	}

	return nil
}
