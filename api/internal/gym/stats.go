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

type ExerciseFailureStats struct {
	ExerciseID      string  `db:"exercise_id" json:"exerciseId"`
	FailureCount    int     `db:"failure_count" json:"failureCount"`
	LastFailureDate *string `db:"last_failure_date" json:"lastFailureDate"`
	TotalSets       int     `db:"total_sets" json:"totalSets"`
}

func GetMusclesFailureStats(db *sqlx.DB, userID uuid.UUID, weeks int) ([]ExerciseFailureStats, error) {
	query := `
		SELECT
			gs.exercise_id,
			COUNT(*) FILTER (WHERE gs.set_type = 'failure') as failure_count,
			MAX(gw.start_time) FILTER (WHERE gs.set_type = 'failure') as last_failure_date,
			COUNT(*) as total_sets
		FROM gym_sets gs
		JOIN gym_workouts gw ON gs.workout_id = gw.id
		WHERE gw.user_id = $1
		AND gw.start_time >= NOW() - ($2 * INTERVAL '1 week')
		GROUP BY gs.exercise_id
		HAVING COUNT(*) FILTER (WHERE gs.set_type = 'failure') > 0
		ORDER BY failure_count DESC`

	var stats []ExerciseFailureStats
	if err := db.Select(&stats, query, userID, weeks); err != nil {
		return nil, fmt.Errorf("muscle failure stats: %w", err)
	}

	if stats == nil {
		stats = []ExerciseFailureStats{}
	}

	return stats, nil
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

type ExerciseSetStats struct {
	Date     string  `db:"date" json:"date"`
	WeightKg float64 `db:"weight_kg" json:"weightKg"`
	Reps     int     `db:"reps" json:"reps"`
	SetType  string  `db:"set_type" json:"setType"`
}

type UserExercise struct {
	ExerciseID string `db:"exercise_id" json:"exerciseId"`
	SetCount   int    `db:"set_count" json:"setCount"`
}

func GetUserExercises(db *sqlx.DB, userID uuid.UUID) ([]UserExercise, error) {
	query := `
		SELECT gs.exercise_id, COUNT(*) as set_count
		FROM gym_sets gs
		JOIN gym_workouts gw ON gs.workout_id = gw.id
		WHERE gw.user_id = $1
		GROUP BY gs.exercise_id
		ORDER BY set_count DESC`

	var exercises []UserExercise
	if err := db.Select(&exercises, query, userID); err != nil {
		return nil, fmt.Errorf("user exercises: %w", err)
	}

	if exercises == nil {
		exercises = []UserExercise{}
	}

	return exercises, nil
}

func GetExerciseStats(db *sqlx.DB, userID uuid.UUID, exerciseID string) ([]ExerciseSetStats, error) {
	query := `
		SELECT
			gw.start_time::date as date,
			gs.weight_kg,
			gs.reps,
			gs.set_type
		FROM gym_sets gs
		JOIN gym_workouts gw ON gs.workout_id = gw.id
		WHERE gw.user_id = $1
		AND gs.exercise_id = $2
		AND gs.weight_kg IS NOT NULL
		AND gs.reps IS NOT NULL
		ORDER BY gw.start_time ASC, gs.position ASC`

	var stats []ExerciseSetStats
	if err := db.Select(&stats, query, userID, exerciseID); err != nil {
		return nil, fmt.Errorf("exercise stats: %w", err)
	}

	if stats == nil {
		stats = []ExerciseSetStats{}
	}

	return stats, nil
}
