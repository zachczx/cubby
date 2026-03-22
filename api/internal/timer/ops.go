package timer

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func GetAllProfiles(db *sqlx.DB, userID uuid.UUID) ([]Profile, error) {
	q := `SELECT * FROM timer_profiles
			WHERE user_id = $1
			ORDER BY created_at ASC`

	var profiles []Profile
	if err := db.Select(&profiles, q, userID); err != nil {
		return nil, fmt.Errorf("get all timer profiles: %w", err)
	}

	if profiles == nil {
		profiles = []Profile{}
	}

	return profiles, nil
}

func NewProfile(db *sqlx.DB, userID uuid.UUID, input ProfileInput) (Profile, error) {
	segJSON, err := json.Marshal(input.Segments)
	if err != nil {
		return Profile{}, fmt.Errorf("marshal segments: %w", err)
	}

	if input.IsDefault {
		if err := clearDefault(db, userID); err != nil {
			return Profile{}, err
		}
	}

	q := `INSERT INTO timer_profiles (user_id, name, segments, is_default)
			VALUES ($1, $2, $3, $4)
			RETURNING id, user_id, name, segments, is_default, created_at, updated_at`

	var p Profile
	err = db.QueryRow(q, userID, input.Name, segJSON, input.IsDefault).Scan(
		&p.ID, &p.UserID, &p.Name, &p.Segments, &p.IsDefault, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return Profile{}, fmt.Errorf("new timer profile: %w", err)
	}

	return p, nil
}

func EditProfile(db *sqlx.DB, userID uuid.UUID, profileID uuid.UUID, input ProfileInput) error {
	segJSON, err := json.Marshal(input.Segments)
	if err != nil {
		return fmt.Errorf("marshal segments: %w", err)
	}

	if input.IsDefault {
		if err := clearDefault(db, userID); err != nil {
			return err
		}
	}

	q := `UPDATE timer_profiles
			SET name = $1, segments = $2, is_default = $3, updated_at = NOW()
			WHERE id = $4 AND user_id = $5`

	if _, err := db.Exec(q, input.Name, segJSON, input.IsDefault, profileID, userID); err != nil {
		return fmt.Errorf("edit timer profile: %w", err)
	}

	return nil
}

func clearDefault(db *sqlx.DB, userID uuid.UUID) error {
	q := `UPDATE timer_profiles SET is_default = FALSE WHERE user_id = $1 AND is_default = TRUE`
	if _, err := db.Exec(q, userID); err != nil {
		return fmt.Errorf("clear default timer profile: %w", err)
	}
	return nil
}

func DeleteProfile(db *sqlx.DB, userID uuid.UUID, profileID uuid.UUID) error {
	q := `DELETE FROM timer_profiles WHERE id = $1 AND user_id = $2`

	if _, err := db.Exec(q, profileID, userID); err != nil {
		return fmt.Errorf("delete timer profile: %w", err)
	}

	return nil
}
