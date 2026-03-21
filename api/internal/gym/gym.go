package gym

import (
	"time"

	"github.com/google/uuid"
)

type Workout struct {
	ID        uuid.UUID  `db:"id"         json:"id"`
	UserID    uuid.UUID  `db:"user_id"    json:"userId"`
	StartTime time.Time  `db:"start_time" json:"startTime"`
	EndTime   *time.Time `db:"end_time"   json:"endTime"`
	Notes     *string    `db:"notes"      json:"notes"`
	CreatedAt time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time  `db:"updated_at" json:"updatedAt"`
	Sets      []Set      `db:"-"          json:"sets"`
}

type WorkoutInput struct {
	StartTime string  `json:"startTime"`
	EndTime   *string `json:"endTime"`
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
	CreatedAt   time.Time `db:"created_at"   json:"createdAt"`
	UpdatedAt   time.Time `db:"updated_at"   json:"updatedAt"`
}

type SetInput struct {
	ExerciseID string   `json:"exerciseId"`
	WeightKg   *float64 `json:"weightKg"`
	Reps       *int16   `json:"reps"`
	SetType    string   `json:"setType"`
}
