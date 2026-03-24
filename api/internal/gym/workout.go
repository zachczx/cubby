package gym

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func NewWorkout(db *sqlx.DB, userID uuid.UUID) (Workout, error) {
	q := `INSERT INTO gym_workouts (user_id)
			VALUES ($1)
			RETURNING id, user_id, start_time, notes, created_at, updated_at`

	var w Workout
	err := db.QueryRow(q, userID).Scan(
		&w.ID, &w.UserID, &w.StartTime, &w.Notes, &w.CreatedAt, &w.UpdatedAt,
	)
	if err != nil {
		return Workout{}, fmt.Errorf("new workout: %w", err)
	}

	w.Sets = []Set{}
	return w, nil
}

func GetAllWorkouts(db *sqlx.DB, userID uuid.UUID) ([]Workout, error) {
	wq := `SELECT * FROM gym_workouts
			WHERE user_id = $1
			ORDER BY start_time DESC`

	var workouts []Workout
	if err := db.Select(&workouts, wq, userID); err != nil {
		return nil, fmt.Errorf("get all workouts: %w", err)
	}

	if len(workouts) == 0 {
		return workouts, nil
	}

	ids := make([]uuid.UUID, len(workouts))
	wMap := make(map[uuid.UUID]int, len(workouts))
	for i, w := range workouts {
		ids[i] = w.ID
		wMap[w.ID] = i
		workouts[i].Sets = []Set{}
	}

	sq := `SELECT * FROM gym_sets
			WHERE workout_id IN (?)
			ORDER BY exercise_id, position ASC`

	query, args, err := sqlx.In(sq, ids)
	if err != nil {
		return nil, fmt.Errorf("get all workouts sets in: %w", err)
	}
	query = db.Rebind(query)

	var sets []Set
	if err := db.Select(&sets, query, args...); err != nil {
		return nil, fmt.Errorf("get all workouts sets: %w", err)
	}

	for _, s := range sets {
		i := wMap[s.WorkoutID]
		workouts[i].Sets = append(workouts[i].Sets, s)
	}

	return workouts, nil
}

func EditWorkout(db *sqlx.DB, userID uuid.UUID, workoutID uuid.UUID, w WorkoutInput) error {
	q := `UPDATE gym_workouts
			SET start_time = $1, notes = $2, updated_at = NOW()
			WHERE id = $3 AND user_id = $4`

	if _, err := db.Exec(q, w.StartTime, w.Notes, workoutID, userID); err != nil {
		return fmt.Errorf("edit workout: %w", err)
	}

	return nil
}

func DeleteWorkout(db *sqlx.DB, userID uuid.UUID, workoutID uuid.UUID) error {
	q := `DELETE FROM gym_workouts WHERE id = $1 AND user_id = $2`

	if _, err := db.Exec(q, workoutID, userID); err != nil {
		return fmt.Errorf("delete workout: %w", err)
	}

	return nil
}
