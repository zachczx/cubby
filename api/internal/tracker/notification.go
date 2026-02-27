package tracker

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/zachczx/cubby/api/internal/notifier"
)

func StartNotifications(db *sqlx.DB, fcm *notifier.FCMClient) error {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		go func() {
			if err := CheckAndNotify(db, fcm); err != nil {
				// return fmt.Errorf("start notif: %w", err)
				log.Println("notification worker error:", err)
			}
		}()
	}

	return nil
}

var ctxTimeout time.Duration = 10

func CheckAndNotify(db *sqlx.DB, fcm *notifier.FCMClient) error {
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
	defer cancel()

	t, err := GetTrackersLast(db)
	if err != nil {
		return fmt.Errorf("get tracker last: %w", err)
	}

	lastDueTrackers, err := CalculateTrackersLastDue(db, t)
	if err != nil {
		return fmt.Errorf("calculateTrackersLastDue: %w", err)
	}

	dueTrackers, err := GetDueTrackerID(lastDueTrackers)
	if err != nil {
		return fmt.Errorf("getDueTrackerID: %w", err)
	}

	userTokens, err := notifier.GetUsersWithTokens(db, dueTrackers)
	if err != nil {
		return fmt.Errorf("getUsersWithTokens: %w", err)
	}

	if err := fcm.SendBatchMessages(db, ctx, userTokens); err != nil {
		return fmt.Errorf("sendBatchMessages: %w", err)
	}

	return nil
}
