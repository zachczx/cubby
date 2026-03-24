package gym

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func GetFavourites(db *sqlx.DB, userID uuid.UUID) ([]string, error) {
	q := `SELECT exercise_id FROM gym_favourite_exercises
			WHERE user_id = $1
			ORDER BY created_at ASC`

	var ids []string
	if err := db.Select(&ids, q, userID); err != nil {
		return nil, fmt.Errorf("get favourites: %w", err)
	}

	return ids, nil
}

func ToggleFavourite(db *sqlx.DB, userID uuid.UUID, exerciseID string) ([]string, error) {
	// Try to delete first; if a row was removed it was already favourited
	delQ := `DELETE FROM gym_favourite_exercises
			WHERE user_id = $1 AND exercise_id = $2`

	result, err := db.Exec(delQ, userID, exerciseID)
	if err != nil {
		return nil, fmt.Errorf("toggle favourite delete: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("toggle favourite rows affected: %w", err)
	}

	if affected == 0 {
		// Was not favourited — insert it
		insQ := `INSERT INTO gym_favourite_exercises (user_id, exercise_id)
				VALUES ($1, $2)`
		if _, err := db.Exec(insQ, userID, exerciseID); err != nil {
			return nil, fmt.Errorf("toggle favourite insert: %w", err)
		}
	}

	return GetFavourites(db, userID)
}
