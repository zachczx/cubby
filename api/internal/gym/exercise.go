package gym

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func NewExercise(db *sqlx.DB, userID uuid.UUID, e ExerciseInput) (Exercise, error) {
	q := `INSERT INTO gym_exercises (name, category, owner_id)
			VALUES ($1, $2, $3)
			RETURNING id, name, category, owner_id, created_at, updated_at`

	var ex Exercise
	err := db.QueryRow(q, e.Name, e.Category, userID).Scan(
		&ex.ID, &ex.Name, &ex.Category, &ex.OwnerID, &ex.CreatedAt, &ex.UpdatedAt,
	)
	if err != nil {
		return Exercise{}, fmt.Errorf("new exercise: %w", err)
	}

	return ex, nil
}

func GetAllExercises(db *sqlx.DB, userID uuid.UUID) ([]Exercise, error) {
	q := `SELECT * FROM gym_exercises
			WHERE owner_id = $1 OR owner_id IS NULL
			ORDER BY name ASC`

	var exercises []Exercise
	if err := db.Select(&exercises, q, userID); err != nil {
		return nil, fmt.Errorf("get all exercises: %w", err)
	}

	return exercises, nil
}

func EditExercise(db *sqlx.DB, userID uuid.UUID, exerciseID uuid.UUID, e ExerciseInput) error {
	q := `UPDATE gym_exercises
			SET name = $1, category = $2, updated_at = NOW()
			WHERE id = $3 AND owner_id = $4`

	if _, err := db.Exec(q, e.Name, e.Category, exerciseID, userID); err != nil {
		return fmt.Errorf("edit exercise: %w", err)
	}

	return nil
}

func DeleteExercise(db *sqlx.DB, userID uuid.UUID, exerciseID uuid.UUID) error {
	q := `DELETE FROM gym_exercises WHERE id = $1 AND owner_id = $2`

	if _, err := db.Exec(q, exerciseID, userID); err != nil {
		return fmt.Errorf("delete exercise: %w", err)
	}

	return nil
}
