package tracker

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/zachczx/cubby/api/internal/user"
)

var trackerDefaults = []TrackerInput{
	{
		Name:         "bedsheet",
		Display:      "Bedsheet",
		Interval:     14,
		IntervalUnit: "day",
		Category:     "household",
		Kind:         "task",
		ActionLabel:  "Changed",
		Pinned:       false,
		Show:         true,
		Icon:         "bed",
	},
	{
		Name:         "petBath",
		Display:      "Dog Bath",
		Interval:     14,
		IntervalUnit: "day",
		Category:     "pet",
		Kind:         "task",
		ActionLabel:  "Bathed",
		Pinned:       false,
		Show:         true,
		Icon:         "shower",
	},
	{
		Name:         "petChewable",
		Display:      "Nexgard",
		Interval:     1,
		IntervalUnit: "month",
		Category:     "pet",
		Kind:         "task",
		ActionLabel:  "Fed",
		Pinned:       false,
		Show:         true,
		Icon:         "shield",
	},
	{
		Name:         "gummy",
		Display:      "Gummy",
		Interval:     2,
		IntervalUnit: "day",
		Category:     "personal",
		Kind:         "task",
		ActionLabel:  "Ate",
		Pinned:       true,
		Show:         true,
		Icon:         "shield",
	},
	{
		Name:         "spray",
		Display:      "Nasal Spray",
		Interval:     3,
		IntervalUnit: "day",
		Category:     "personal",
		Kind:         "task",
		ActionLabel:  "Sprayed",
		Pinned:       true,
		Show:         true,
		Icon:         "bottle",
	},
	{
		Name:         "towel",
		Display:      "Towel Wash",
		Interval:     5,
		IntervalUnit: "day",
		Category:     "household",
		Kind:         "task",
		ActionLabel:  "Washed",
		Pinned:       true,
		Show:         true,
		Icon:         "washer",
	},
}

type DefaultService struct{}

func (DefaultService) CreateDefaults(db *sqlx.DB, userID uuid.UUID) error {
	familyID, err := user.GetUserFamilyID(db, userID)
	if err != nil {
		return fmt.Errorf("getting user family id: %w", err)
	}

	for _, d := range trackerDefaults {
		t := Tracker{
			User:         userID,
			Family:       familyID,
			Name:         d.Name,
			Display:      d.Display,
			Interval:     d.Interval,
			IntervalUnit: d.IntervalUnit,
			Category:     d.Category,
			Kind:         d.Kind,
			ActionLabel:  d.ActionLabel,
			Pinned:       d.Pinned,
			Show:         d.Show,
			Icon:         d.Icon,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		if _, err := New(db, t); err != nil {
			return fmt.Errorf("creating default tracker %s: %w", d.Name, err)
		}
	}

	return nil
}
