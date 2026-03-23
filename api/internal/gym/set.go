package gym

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func NewSet(db *sqlx.DB, userID uuid.UUID, workoutID uuid.UUID, s SetInput) (Set, error) {
	q := `INSERT INTO gym_sets (workout_id, exercise_id, weight_kg, reps, set_type, position)
			SELECT $1, $2, $3, $4, $5,
				COALESCE((SELECT MAX(position) + 1 FROM gym_sets WHERE workout_id = $1 AND exercise_id = $2), 0)
			FROM gym_workouts
			WHERE id = $1 AND user_id = $6
			RETURNING id, workout_id, exercise_id, weight_kg, reps, set_type, is_completed, position, created_at, updated_at`

	var set Set
	err := db.QueryRow(q, workoutID, s.ExerciseID, s.WeightKg, s.Reps, s.SetType, userID).Scan(
		&set.ID, &set.WorkoutID, &set.ExerciseID, &set.WeightKg, &set.Reps,
		&set.SetType, &set.IsCompleted, &set.Position, &set.CreatedAt, &set.UpdatedAt,
	)
	if err != nil {
		return Set{}, fmt.Errorf("new set: %w", err)
	}

	return set, nil
}

func EditSet(db *sqlx.DB, userID uuid.UUID, setID uuid.UUID, s SetInput) error {
	q := `UPDATE gym_sets
			SET exercise_id = $1, weight_kg = $2, reps = $3, set_type = $4, updated_at = NOW()
			FROM gym_workouts
			WHERE gym_sets.id = $5
			AND gym_sets.workout_id = gym_workouts.id
			AND gym_workouts.user_id = $6`

	if _, err := db.Exec(q, s.ExerciseID, s.WeightKg, s.Reps, s.SetType, setID, userID); err != nil {
		return fmt.Errorf("edit set: %w", err)
	}

	return nil
}

func ReorderSet(db *sqlx.DB, userID uuid.UUID, input ReorderSetInput) error {
	// Get the set and its neighbor in the given direction
	var current Set
	q := `SELECT gs.* FROM gym_sets gs
			JOIN gym_workouts gw ON gs.workout_id = gw.id
			WHERE gs.id = $1 AND gw.user_id = $2`
	if err := db.Get(&current, q, input.SetID, userID); err != nil {
		return fmt.Errorf("reorder set get current: %w", err)
	}

	var neighbor Set
	var nq string
	if input.Direction == "up" {
		nq = `SELECT * FROM gym_sets
				WHERE workout_id = $1 AND exercise_id = $2 AND position < $3
				ORDER BY position DESC LIMIT 1`
	} else {
		nq = `SELECT * FROM gym_sets
				WHERE workout_id = $1 AND exercise_id = $2 AND position > $3
				ORDER BY position ASC LIMIT 1`
	}
	if err := db.Get(&neighbor, nq, current.WorkoutID, current.ExerciseID, current.Position); err != nil {
		return fmt.Errorf("reorder set get neighbor: %w", err)
	}

	// Swap positions
	swapQ := `UPDATE gym_sets SET position = CASE
				WHEN id = $1 THEN $2
				WHEN id = $3 THEN $4
			END, updated_at = NOW()
			WHERE id IN ($1, $3)`
	if _, err := db.Exec(swapQ, current.ID, neighbor.Position, neighbor.ID, current.Position); err != nil {
		return fmt.Errorf("reorder set swap: %w", err)
	}

	return nil
}

func DeleteSet(db *sqlx.DB, userID uuid.UUID, setID uuid.UUID) error {
	q := `DELETE FROM gym_sets
			USING gym_workouts
			WHERE gym_sets.id = $1
			AND gym_sets.workout_id = gym_workouts.id
			AND gym_workouts.user_id = $2`

	if _, err := db.Exec(q, setID, userID); err != nil {
		return fmt.Errorf("delete set: %w", err)
	}

	return nil
}
