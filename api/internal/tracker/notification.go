package tracker

import (
	"fmt"
	"log"
	"time"

	"github.com/zachczx/cubby/api/internal/server"
)

func StartNotifications(s *server.Service) error {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		go func() {
			if err := CheckAndNotify(s); err != nil {
				// return fmt.Errorf("start notif: %w", err)
				log.Println("notification worker error:", err)
			}
		}()
	}

	return nil
}

func CheckAndNotify(s *server.Service) error {
	t, err := GetTrackersLast(s.DB)
	if err != nil {
		return fmt.Errorf("get tracker last: %w", err)
	}

	newT, err := CalculateTrackersLastDue(s.DB, t)
	if err != nil {
		return fmt.Errorf("calculateTrackersLastDue: %w", err)
	}

	fmt.Println(newT)

	return nil
}
