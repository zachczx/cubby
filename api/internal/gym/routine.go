package gym

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func GetAllRoutines(db *sqlx.DB, userID uuid.UUID) ([]Routine, error) {
	rq := `SELECT * FROM gym_routines
			WHERE user_id = $1
			ORDER BY position ASC, created_at ASC`

	var routines []Routine
	if err := db.Select(&routines, rq, userID); err != nil {
		return nil, fmt.Errorf("get all routines: %w", err)
	}

	if len(routines) == 0 {
		return routines, nil
	}

	ids := make([]uuid.UUID, len(routines))
	rMap := make(map[uuid.UUID]int, len(routines))
	for i, r := range routines {
		ids[i] = r.ID
		rMap[r.ID] = i
		routines[i].Exercises = []RoutineExercise{}
	}

	eq := `SELECT * FROM gym_routine_exercises
			WHERE routine_id IN (?)
			ORDER BY position ASC`

	query, args, err := sqlx.In(eq, ids)
	if err != nil {
		return nil, fmt.Errorf("get all routines exercises in: %w", err)
	}
	query = db.Rebind(query)

	var exercises []RoutineExercise
	if err := db.Select(&exercises, query, args...); err != nil {
		return nil, fmt.Errorf("get all routines exercises: %w", err)
	}

	for _, e := range exercises {
		i := rMap[e.RoutineID]
		routines[i].Exercises = append(routines[i].Exercises, e)
	}

	return routines, nil
}

func NewRoutine(db *sqlx.DB, userID uuid.UUID, input RoutineInput) (Routine, error) {
	q := `INSERT INTO gym_routines (user_id, name, position)
			VALUES ($1, $2, COALESCE((SELECT MAX(position) + 1 FROM gym_routines WHERE user_id = $1), 0))
			RETURNING id, user_id, name, position, created_at, updated_at`

	var r Routine
	err := db.QueryRow(q, userID, input.Name).Scan(
		&r.ID, &r.UserID, &r.Name, &r.Position, &r.CreatedAt, &r.UpdatedAt,
	)
	if err != nil {
		return Routine{}, fmt.Errorf("new routine: %w", err)
	}

	r.Exercises = []RoutineExercise{}
	return r, nil
}

func EditRoutine(db *sqlx.DB, userID uuid.UUID, routineID uuid.UUID, input RoutineInput) error {
	q := `UPDATE gym_routines
			SET name = $1, updated_at = NOW()
			WHERE id = $2 AND user_id = $3`

	if _, err := db.Exec(q, input.Name, routineID, userID); err != nil {
		return fmt.Errorf("edit routine: %w", err)
	}

	return nil
}

func DeleteRoutine(db *sqlx.DB, userID uuid.UUID, routineID uuid.UUID) error {
	q := `DELETE FROM gym_routines WHERE id = $1 AND user_id = $2`

	if _, err := db.Exec(q, routineID, userID); err != nil {
		return fmt.Errorf("delete routine: %w", err)
	}

	return nil
}

func AddRoutineExercise(db *sqlx.DB, userID uuid.UUID, routineID uuid.UUID, input RoutineExerciseInput) (RoutineExercise, error) {
	q := `INSERT INTO gym_routine_exercises (routine_id, exercise_id, sets, position)
			SELECT $1, $2, $3,
				COALESCE((SELECT MAX(position) + 1 FROM gym_routine_exercises WHERE routine_id = $1), 0)
			FROM gym_routines
			WHERE id = $1 AND user_id = $4
			RETURNING id, routine_id, exercise_id, sets, position, created_at, updated_at`

	var e RoutineExercise
	err := db.QueryRow(q, routineID, input.ExerciseID, input.Sets, userID).Scan(
		&e.ID, &e.RoutineID, &e.ExerciseID, &e.Sets, &e.Position, &e.CreatedAt, &e.UpdatedAt,
	)
	if err != nil {
		return RoutineExercise{}, fmt.Errorf("add routine exercise: %w", err)
	}

	return e, nil
}

func EditRoutineExercise(db *sqlx.DB, userID uuid.UUID, exerciseID uuid.UUID, input RoutineExerciseInput) error {
	q := `UPDATE gym_routine_exercises
			SET sets = $1, updated_at = NOW()
			FROM gym_routines
			WHERE gym_routine_exercises.id = $2
			AND gym_routine_exercises.routine_id = gym_routines.id
			AND gym_routines.user_id = $3`

	if _, err := db.Exec(q, input.Sets, exerciseID, userID); err != nil {
		return fmt.Errorf("edit routine exercise: %w", err)
	}

	return nil
}

func RemoveRoutineExercise(db *sqlx.DB, userID uuid.UUID, exerciseID uuid.UUID) error {
	q := `DELETE FROM gym_routine_exercises
			USING gym_routines
			WHERE gym_routine_exercises.id = $1
			AND gym_routine_exercises.routine_id = gym_routines.id
			AND gym_routines.user_id = $2`

	if _, err := db.Exec(q, exerciseID, userID); err != nil {
		return fmt.Errorf("remove routine exercise: %w", err)
	}

	return nil
}

func ReorderRoutineExercise(db *sqlx.DB, userID uuid.UUID, input ReorderRoutineExerciseInput) error {
	var current RoutineExercise
	q := `SELECT gre.* FROM gym_routine_exercises gre
			JOIN gym_routines gr ON gre.routine_id = gr.id
			WHERE gre.id = $1 AND gr.user_id = $2`
	if err := db.Get(&current, q, input.ExerciseID, userID); err != nil {
		return fmt.Errorf("reorder routine exercise get current: %w", err)
	}

	var neighbor RoutineExercise
	var nq string
	if input.Direction == "up" {
		nq = `SELECT * FROM gym_routine_exercises
				WHERE routine_id = $1 AND position < $2
				ORDER BY position DESC LIMIT 1`
	} else {
		nq = `SELECT * FROM gym_routine_exercises
				WHERE routine_id = $1 AND position > $2
				ORDER BY position ASC LIMIT 1`
	}
	if err := db.Get(&neighbor, nq, current.RoutineID, current.Position); err != nil {
		return fmt.Errorf("reorder routine exercise get neighbor: %w", err)
	}

	swapQ := `UPDATE gym_routine_exercises SET position = $1::smallint, updated_at = NOW() WHERE id = $2`
	if _, err := db.Exec(swapQ, neighbor.Position, current.ID); err != nil {
		return fmt.Errorf("reorder routine exercise swap current: %w", err)
	}
	if _, err := db.Exec(swapQ, current.Position, neighbor.ID); err != nil {
		return fmt.Errorf("reorder routine exercise swap neighbor: %w", err)
	}

	return nil
}

func StartWorkoutFromRoutine(db *sqlx.DB, userID uuid.UUID, routineID uuid.UUID) (Workout, error) {
	tx, err := db.Beginx()
	if err != nil {
		return Workout{}, fmt.Errorf("start workout from routine begin tx: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	// Verify routine belongs to user and get exercises
	var exercises []RoutineExercise
	eq := `SELECT gre.* FROM gym_routine_exercises gre
			JOIN gym_routines gr ON gre.routine_id = gr.id
			WHERE gr.id = $1 AND gr.user_id = $2
			ORDER BY gre.position ASC`
	if err := tx.Select(&exercises, eq, routineID, userID); err != nil {
		return Workout{}, fmt.Errorf("start workout from routine get exercises: %w", err)
	}

	// Create the workout
	wq := `INSERT INTO gym_workouts (user_id)
			VALUES ($1)
			RETURNING id, user_id, start_time, notes, created_at, updated_at`
	var w Workout
	if err := tx.QueryRow(wq, userID).Scan(
		&w.ID, &w.UserID, &w.StartTime, &w.Notes, &w.CreatedAt, &w.UpdatedAt,
	); err != nil {
		return Workout{}, fmt.Errorf("start workout from routine create workout: %w", err)
	}

	// Get last-used weight/reps for each exercise
	exerciseIDs := make([]string, len(exercises))
	for i, e := range exercises {
		exerciseIDs[i] = e.ExerciseID
	}

	type lastUsed struct {
		ExerciseID string   `db:"exercise_id"`
		WeightKg   *float64 `db:"weight_kg"`
		Reps       *int16   `db:"reps"`
		SetType    string   `db:"set_type"`
	}

	lastUsedMap := make(map[string]lastUsed)
	if len(exerciseIDs) > 0 {
		lq := `SELECT DISTINCT ON (gs.exercise_id) gs.exercise_id, gs.weight_kg, gs.reps, gs.set_type
				FROM gym_sets gs
				JOIN gym_workouts gw ON gs.workout_id = gw.id
				WHERE gw.user_id = ? AND gs.exercise_id IN (?)
				ORDER BY gs.exercise_id, gw.start_time DESC, gs.position DESC`

		allArgs := []any{userID}
		allArgs = append(allArgs, toAnySlice(exerciseIDs)...)

		query, args, err := sqlx.In(lq, allArgs...)
		if err != nil {
			return Workout{}, fmt.Errorf("start workout from routine last used in: %w", err)
		}
		query = db.Rebind(query)

		var rows []lastUsed
		if err := tx.Select(&rows, query, args...); err != nil {
			return Workout{}, fmt.Errorf("start workout from routine last used: %w", err)
		}
		for _, r := range rows {
			lastUsedMap[r.ExerciseID] = r
		}
	}

	// Insert sets for each routine exercise
	setQ := `INSERT INTO gym_sets (workout_id, exercise_id, weight_kg, reps, set_type, position)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id, workout_id, exercise_id, weight_kg, reps, set_type, is_completed, position, created_at, updated_at`

	w.Sets = []Set{}
	pos := int16(0)
	for _, re := range exercises {
		lu := lastUsedMap[re.ExerciseID]
		setType := "working"
		if lu.SetType != "" {
			setType = lu.SetType
		}

		for range re.Sets {
			var s Set
			if err := tx.QueryRow(setQ, w.ID, re.ExerciseID, lu.WeightKg, lu.Reps, setType, pos).Scan(
				&s.ID, &s.WorkoutID, &s.ExerciseID, &s.WeightKg, &s.Reps,
				&s.SetType, &s.IsCompleted, &s.Position, &s.CreatedAt, &s.UpdatedAt,
			); err != nil {
				return Workout{}, fmt.Errorf("start workout from routine insert set: %w", err)
			}
			w.Sets = append(w.Sets, s)
			pos++
		}
	}

	if err := tx.Commit(); err != nil {
		return Workout{}, fmt.Errorf("start workout from routine commit: %w", err)
	}

	return w, nil
}

func toAnySlice(ss []string) []any {
	out := make([]any, len(ss))
	for i, s := range ss {
		out[i] = s
	}
	return out
}
