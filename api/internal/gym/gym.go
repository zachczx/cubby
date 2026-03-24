package gym

import (
	"time"

	"github.com/google/uuid"
)

type Workout struct {
	ID        uuid.UUID `db:"id"         json:"id"`
	UserID    uuid.UUID `db:"user_id"    json:"userId"`
	StartTime time.Time `db:"start_time" json:"startTime"`
	Notes     *string   `db:"notes"      json:"notes"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
	Sets      []Set     `db:"-"          json:"sets"`
}

type WorkoutInput struct {
	StartTime string  `json:"startTime"`
	Notes     *string `json:"notes"`
}

type Set struct {
	ID          uuid.UUID `db:"id"           json:"id"`
	WorkoutID   uuid.UUID `db:"workout_id"   json:"workoutId"`
	ExerciseID  string    `db:"exercise_id"  json:"exerciseId"`
	WeightKg    *float64  `db:"weight_kg"    json:"weightKg"`
	Reps        *int16    `db:"reps"         json:"reps"`
	SetType     string    `db:"set_type"     json:"setType"`
	IsCompleted bool      `db:"is_completed" json:"isCompleted"`
	Position    int16     `db:"position"     json:"position"`
	CreatedAt   time.Time `db:"created_at"   json:"createdAt"`
	UpdatedAt   time.Time `db:"updated_at"   json:"updatedAt"`
}

type ReorderSetInput struct {
	SetID     uuid.UUID `json:"setId"`
	Direction string    `json:"direction"`
}

type SetInput struct {
	ExerciseID string   `json:"exerciseId"`
	WeightKg   *float64 `json:"weightKg"`
	Reps       *int16   `json:"reps"`
	SetType    string   `json:"setType"`
}

type FavouriteExercise struct {
	UserID     uuid.UUID `db:"user_id"`
	ExerciseID string    `db:"exercise_id"`
	CreatedAt  time.Time `db:"created_at"`
}
