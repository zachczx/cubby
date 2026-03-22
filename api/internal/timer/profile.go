package timer

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Segment struct {
	Label          string `json:"label"`
	DefaultSeconds int    `json:"defaultSeconds"`
}

type Profile struct {
	ID        uuid.UUID       `db:"id"         json:"id"`
	UserID    uuid.UUID       `db:"user_id"    json:"userId"`
	Name      string          `db:"name"       json:"name"`
	Segments  json.RawMessage `db:"segments"   json:"segments"`
	IsDefault bool            `db:"is_default" json:"isDefault"`
	CreatedAt time.Time       `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time       `db:"updated_at" json:"updatedAt"`
}

type ProfileInput struct {
	Name      string    `json:"name"`
	Segments  []Segment `json:"segments"`
	IsDefault bool      `json:"isDefault"`
}
