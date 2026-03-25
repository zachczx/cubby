package gym

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type WorkoutSummary struct {
	TotalWorkoutsThisMonth int             `json:"totalWorkoutsThisMonth"`
	TotalVolumeThisMonth   float64         `json:"totalVolumeThisMonth"`
	TotalSetsThisMonth     int             `json:"totalSetsThisMonth"`
	TopExercises           []ExerciseCount `json:"topExercises"`
	FailureExercises       []ExerciseCount `json:"failureExercises"`
}

type ExerciseCount struct {
	ExerciseID string `db:"exercise_id" json:"exerciseId"`
	Count      int    `db:"count"        json:"count"`
}

func GetSummary(db *sqlx.DB, userID uuid.UUID) (WorkoutSummary, error) {
	var summary WorkoutSummary

	workoutsQ := `SELECT COUNT(*) FROM gym_workouts
			WHERE user_id = $1
			AND start_time >= date_trunc('month', NOW())`
	if err := db.Get(&summary.TotalWorkoutsThisMonth, workoutsQ, userID); err != nil {
		return summary, fmt.Errorf("summary workouts count: %w", err)
	}

	volumeQ := `SELECT COALESCE(SUM(gs.weight_kg * gs.reps), 0), COUNT(*)
			FROM gym_sets gs
			JOIN gym_workouts gw ON gs.workout_id = gw.id
			WHERE gw.user_id = $1
			AND gw.start_time >= date_trunc('month', NOW())
			AND gs.weight_kg IS NOT NULL
			AND gs.reps IS NOT NULL`
	if err := db.QueryRow(volumeQ, userID).Scan(&summary.TotalVolumeThisMonth, &summary.TotalSetsThisMonth); err != nil {
		return summary, fmt.Errorf("summary volume: %w", err)
	}

	topExercisesQ := `SELECT gs.exercise_id, COUNT(*) as count
			FROM gym_sets gs
			JOIN gym_workouts gw ON gs.workout_id = gw.id
			WHERE gw.user_id = $1
			GROUP BY gs.exercise_id
			ORDER BY count DESC
			LIMIT 5`
	if err := db.Select(&summary.TopExercises, topExercisesQ, userID); err != nil {
		return summary, fmt.Errorf("summary top exercises: %w", err)
	}

	failureQ := `SELECT gs.exercise_id, COUNT(*) as count
			FROM gym_sets gs
			JOIN gym_workouts gw ON gs.workout_id = gw.id
			WHERE gw.user_id = $1
			AND gs.set_type = 'failure'
			GROUP BY gs.exercise_id
			ORDER BY count DESC
			LIMIT 10`
	if err := db.Select(&summary.FailureExercises, failureQ, userID); err != nil {
		return summary, fmt.Errorf("summary failure exercises: %w", err)
	}

	return summary, nil
}

type WorkoutCalendarEntry struct {
	WorkoutID     string         `db:"workout_id"      json:"workoutId"`
	StartTime     string         `db:"start_time"      json:"startTime"`
	ExerciseCount int            `db:"exercise_count"  json:"exerciseCount"`
	SetCount      int            `db:"set_count"       json:"setCount"`
	ExerciseIDs   pq.StringArray `db:"exercise_ids"    json:"exerciseIds"`
}

func GetCalendarWorkouts(db *sqlx.DB, userID uuid.UUID) ([]WorkoutCalendarEntry, error) {
	// FILTER excludes null values from LEFT JOIN when workout doesnt have any sets
	calendarQ := `SELECT
			gw.id as workout_id,
			gw.start_time,
			COUNT(DISTINCT gs.exercise_id) as exercise_count,
			COUNT(gs.id) as set_count,
			COALESCE(ARRAY_AGG(DISTINCT gs.exercise_id) FILTER (WHERE gs.exercise_id IS NOT NULL), '{}') as exercise_ids
		FROM gym_workouts gw
		LEFT JOIN gym_sets gs ON gs.workout_id = gw.id
		WHERE gw.user_id = $1
		GROUP BY gw.id, gw.start_time
		ORDER BY gw.start_time DESC`

	var entries []WorkoutCalendarEntry
	if err := db.Select(&entries, calendarQ, userID); err != nil {
		return nil, fmt.Errorf("calendar workouts: %w", err)
	}

	if entries == nil {
		entries = []WorkoutCalendarEntry{}
	}

	return entries, nil
}
