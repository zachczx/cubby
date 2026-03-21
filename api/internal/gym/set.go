package gym

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func NewSet(db *sqlx.DB, userID uuid.UUID, workoutID uuid.UUID, s SetInput) (Set, error) {
	q := `INSERT INTO gym_sets (workout_id, exercise_id, weight_kg, reps, set_type)
			SELECT $1, $2, $3, $4, $5
			FROM gym_workouts
			WHERE id = $1 AND user_id = $6
			RETURNING id, workout_id, exercise_id, weight_kg, reps, set_type, is_completed, created_at, updated_at`

	var set Set
	err := db.QueryRow(q, workoutID, s.ExerciseID, s.WeightKg, s.Reps, s.SetType, userID).Scan(
		&set.ID, &set.WorkoutID, &set.ExerciseID, &set.WeightKg, &set.Reps,
		&set.SetType, &set.IsCompleted, &set.CreatedAt, &set.UpdatedAt,
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
