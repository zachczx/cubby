package timer

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

var profileDefaults = []ProfileInput{
	{
		Name:      "TV Time",
		IsDefault: true,
		Segments: []Segment{
			{Label: "TV", DefaultSeconds: 1800},
			{Label: "Rest", DefaultSeconds: 300},
			{Label: "TV", DefaultSeconds: 1800},
		},
	},
}

func CreateDefaults(db *sqlx.DB, userID uuid.UUID) error {
	for _, d := range profileDefaults {
		if _, err := NewProfile(db, userID, d); err != nil {
			return fmt.Errorf("creating default timer profile %s: %w", d.Name, err)
		}
	}

	return nil
}
