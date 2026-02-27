package tracker

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Tracker struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	Owner        uuid.UUID  `json:"-" db:"owner_id"`
	Family       uuid.UUID  `json:"familyId" db:"family_id"`
	Name         string     `json:"name" db:"name"`
	Display      string     `json:"display" db:"display"`
	Interval     int        `json:"interval" db:"interval"`
	IntervalUnit string     `json:"intervalUnit" db:"interval_unit"`
	Category     string     `json:"category" db:"category"`
	Kind         string     `json:"kind" db:"kind"`
	ActionLabel  string     `json:"actionLabel" db:"action_label"`
	Pinned       bool       `json:"pinned" db:"pinned"`
	Show         bool       `json:"show" db:"show"`
	Icon         string     `json:"icon" db:"icon"`
	StartDate    *time.Time `json:"startDate,omitempty" db:"start_date"`
	Cost         *float64   `json:"cost,omitempty" db:"cost"`
	CreatedAt    time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt    time.Time  `json:"updatedAt" db:"updated_at"`

	FamilyName string `json:"familyName" db:"family_name"`
	IsOwner    bool   `json:"isOwner" db:"-"`
}

type TrackerInput struct {
	Name         string   `json:"name"`
	Display      string   `json:"display"`
	Interval     int      `json:"interval"`
	IntervalUnit string   `json:"intervalUnit"`
	Category     string   `json:"category"`
	Kind         string   `json:"kind"`
	ActionLabel  string   `json:"actionLabel"`
	Pinned       bool     `json:"pinned"`
	Show         bool     `json:"show"`
	Icon         string   `json:"icon"`
	StartDate    string   `json:"startDate"`
	Cost         *float64 `json:"cost"`
}

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

type TrackerLatestEntry struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	Owner        uuid.UUID  `json:"-" db:"owner_id"`
	Family       uuid.UUID  `json:"familyId" db:"family_id"`
	Name         string     `json:"name" db:"name"`
	Display      string     `json:"display" db:"display"`
	Interval     int        `json:"interval" db:"interval"`
	IntervalUnit string     `json:"intervalUnit" db:"interval_unit"`
	Category     string     `json:"category" db:"category"`
	Kind         string     `json:"kind" db:"kind"`
	ActionLabel  string     `json:"actionLabel" db:"action_label"`
	Pinned       bool       `json:"pinned" db:"pinned"`
	Show         bool       `json:"show" db:"show"`
	Icon         string     `json:"icon" db:"icon"`
	StartDate    *time.Time `json:"startDate,omitempty" db:"start_date"`
	Cost         *float64   `json:"cost,omitempty" db:"cost"`
	CreatedAt    time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt    time.Time  `json:"updatedAt" db:"updated_at"`

	FamilyName string `json:"familyName" db:"family_name"`
	IsOwner    bool   `json:"isOwner" db:"-"`

	LastEntry        *time.Time `db:"last_entry" json:"lastEntry"`
	LastInterval     *int       `json:"lastInterval" db:"last_interval"`
	LastIntervalUnit *string    `json:"lastIntervalUnit" db:"last_interval_unit"`
	DueStatus        *string    `json:"dueStatus" db:"-"`
}

func GetTrackersLast(db *sqlx.DB) ([]TrackerLatestEntry, error) {
	q := `SELECT t.*, e.last_entry, e.last_interval, e.last_interval_unit, f.name AS family_name FROM trackers t
			JOIN (
				SELECT tracker_id, MAX(performed_at) AS last_entry, interval AS last_interval, interval_unit AS last_interval_unit FROM entries GROUP BY tracker_id, last_interval, last_interval_unit
			) AS e ON t.id = e.tracker_id
			JOIN families f ON t.family_id = f.id`

	var t []TrackerLatestEntry

	if err := db.Select(&t, q); err != nil {
		return nil, fmt.Errorf("get trackers list: %w", err)
	}

	return t, nil
}

func CalculateTrackersLastDue(db *sqlx.DB, tDB []TrackerLatestEntry) ([]TrackerLatestEntry, error) {
	var newT []TrackerLatestEntry

	newT = tDB

	for i := range tDB {
		var threshold time.Time

		if tDB[i].LastEntry == nil || tDB[i].LastInterval == nil || tDB[i].LastIntervalUnit == nil {
			continue
		}

		itv := &tDB[i].Interval

		switch tDB[i].IntervalUnit {
		case "day":
			threshold = tDB[i].LastEntry.Add(time.Duration(*itv) * 24 * time.Hour)
			threshold = threshold.Add(time.Hour * 6)

		case "month":
			threshold = tDB[i].LastEntry.AddDate(0, int(*itv), 0)
			threshold = threshold.Add(time.Hour * 12)

		case "year":
			threshold = tDB[i].LastEntry.AddDate(int(*itv), 0, 0)
			threshold = threshold.Add(time.Hour * 24)
		}

		if time.Now().After(threshold) {
			new := "due"
			newT[i].DueStatus = &new
		} else {
			new := "ok"
			newT[i].DueStatus = &new
		}
	}

	return newT, nil
}

func GetDueTrackerID(trackers []TrackerLatestEntry) ([]uuid.UUID, error) {
	var due []uuid.UUID

	for _, t := range trackers {
		if *t.DueStatus == "due" {
			due = append(due, t.ID)
		}
	}
	return due, nil
}

