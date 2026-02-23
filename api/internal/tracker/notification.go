package tracker

import (
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

func StartNotifications(db *sqlx.DB) error {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		go func() {
			if err := CheckAndNotify(db); err != nil {
				// return fmt.Errorf("start notif: %w", err)
				log.Println("notification worker error:", err)
			}
		}()
	}

	return nil
}

func CheckAndNotify(db *sqlx.DB) error {
	t, err := GetTrackersLast(db)
	if err != nil {
		return fmt.Errorf("get tracker last: %w", err)
	}

	newT, err := CalculateTrackersLastDue(db, t)
	if err != nil {
		return fmt.Errorf("calculateTrackersLastDue: %w", err)
	}

	_ = newT

	return nil
}
