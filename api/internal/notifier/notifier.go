package notifier

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/auth/credentials"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"google.golang.org/api/option"
)

type PushToken struct {
	ID        string    `db:"id"`
	UserID    string    `db:"user_id"`
	Token     string    `db:"token"`
	Platform  *string   `db:"platform"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type UserToken struct {
	UserID             string    `db:"user_id"`
	UserName           string    `db:"user_name"`
	Token              string    `db:"token"`
	Platform           *string   `db:"platform"`
	TrackerID          uuid.UUID `db:"tracker_id"`
	TrackerDisplayName string    `db:"tracker_display"`
}

type DueTrackerNotification struct {
	Token       string    `db:"token"`
	UserID      uuid.UUID `db:"user_id"`
	TrackerName string    `db:"tracker_name"`
	DueDate     time.Time `db:"due_date"`
}

type UserNotificationData struct {
	UserID      uuid.UUID
	Tokens      []string
	DueTrackers []string
}

type FCMClient struct {
	client *messaging.Client
}

func NewFCMClient(ctx context.Context) (*FCMClient, error) {
	creds, err := credentials.DetectDefault(&credentials.DetectOptions{
		CredentialsJSON: []byte(os.Getenv("FIREBASE_CREDENTIALS_JSON")),
		Scopes: []string{
			"https://www.googleapis.com/auth/firebase.messaging",
			"https://www.googleapis.com/auth/cloud-platform",
		},
	})
	if err != nil {
		return nil, err
	}

	app, err := firebase.NewApp(ctx, &firebase.Config{
		ProjectID: "cubbydotdev",
	}, option.WithAuthCredentials(creds))
	if err != nil {
		return nil, err
	}

	client, err := app.Messaging(ctx)
	if err != nil {
		return nil, err
	}

	return &FCMClient{client: client}, nil
}

var notificationWindowHours time.Duration = -6

func GetUsersWithTokens(db *sqlx.DB, trackerIDs []uuid.UUID) ([]UserToken, error) {
	var tokens []UserToken
	notificationThreshold := time.Now().Add(notificationWindowHours * time.Hour)

	q := `SELECT 
				pt.token,
				pt.platform,
				u.id AS user_id,
				COALESCE(u.name, '') AS user_name,
				t.display AS tracker_display,
				t.id AS tracker_id 
			FROM trackers t
			LEFT JOIN 
				(
				SELECT id AS family_id, owner_id AS user_id FROM families
				UNION
				SELECT family_id, user_id FROM families_users
				) AS fu ON t.family_id = fu.family_id
			LEFT JOIN notification_logs nl ON t.id = nl.tracker_id AND fu.user_id = nl.user_id
			LEFT JOIN users u ON fu.user_id = u.id
			LEFT JOIN push_tokens pt ON u.id = pt.user_id
			WHERE t.id IN (?) 
			AND (nl.id IS NULL OR nl.updated_at < (?))`

	query, args, err := sqlx.In(q, trackerIDs, notificationThreshold)
	if err != nil {
		return nil, fmt.Errorf("getUsersWithTokens In: %w", err)
	}

	query = db.Rebind(query)

	if err := db.Select(&tokens, query, args...); err != nil {
		return nil, fmt.Errorf("getUsersWithTokens select: %w", err)
	}

	return tokens, nil
}

type NotificationMessage struct {
	TrackerID          []uuid.UUID
	TrackerDisplayName []string
}

/*
Using a map lets me send only 1 message for each token/device. Don't want the same device receiving multiple alerts.
*/
func BatchMessageBuilder(userTokens []UserToken) ([]*messaging.Message, error) {
	msges := make(map[string]NotificationMessage)

	for _, u := range userTokens {
		m, exists := msges[u.Token]
		if !exists {
			msges[u.Token] = NotificationMessage{TrackerID: []uuid.UUID{}, TrackerDisplayName: []string{}}
		}

		m.TrackerDisplayName = append(m.TrackerDisplayName, u.TrackerDisplayName)
		m.TrackerID = append(m.TrackerID, u.TrackerID)

		msges[u.Token] = m
	}

	var FCMMessages []*messaging.Message

	for t, m := range msges {
		names := strings.Join(m.TrackerDisplayName, ", ")

		FCMMessages = append(FCMMessages, &messaging.Message{
			Token: t,
			Notification: &messaging.Notification{
				Title: "Cubby Reminder",
				Body:  fmt.Sprintf("Trackers due: %s", names),
			},
			Android: &messaging.AndroidConfig{
				Priority: "high",
			},
			APNS: &messaging.APNSConfig{
				Payload: &messaging.APNSPayload{
					Aps: &messaging.Aps{
						Sound: "default",
					},
				},
			},
		})
	}

	return FCMMessages, nil
}

func (f *FCMClient) SendBatchMessages(db *sqlx.DB, ctx context.Context, userTokens []UserToken) error {
	if len(userTokens) == 0 {
		return nil
	}

	messages, err := BatchMessageBuilder(userTokens)
	if err != nil {
		return err
	}

	if len(messages) == 0 {
		return nil
	}

	batchResponse, err := f.client.SendEach(ctx, messages)
	if err != nil {
		return err
	}
	if batchResponse.FailureCount > 0 {
		for _, resp := range batchResponse.Responses {
			if !resp.Success {
				fmt.Printf("batch msg err: %v\r\n", resp.Error)
			}
		}
		return fmt.Errorf("send batch: %v failures out of %v", batchResponse.FailureCount, len(messages))
	}

	if err := UpdateNotificationLogs(db, userTokens); err != nil {
		return fmt.Errorf("send batch: %w", err)
	}

	return nil
}

func UpdateNotificationLogs(db *sqlx.DB, userTokens []UserToken) error {
	q := `INSERT INTO notification_logs (tracker_id, user_id) 
			VALUES ($1, $2)
			ON CONFLICT (tracker_id, user_id) DO UPDATE 
			SET updated_at = NOW()`

	for _, ut := range userTokens {
		if _, err := db.Exec(q, ut.TrackerID, ut.UserID); err != nil {
			return fmt.Errorf("updateNotificationLogs (tracker id: %v): %w", ut.TrackerID, err)
		}
	}

	return nil
}

func SavePushToken(db *sqlx.DB, userID uuid.UUID, token string, platform string) error {
	tx, err := db.Beginx()
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	qDel := `DELETE FROM push_tokens WHERE user_id = $1 AND platform = $2`
	if _, err := tx.Exec(qDel, userID, platform); err != nil {
		return fmt.Errorf("del stale tokens: %w", err)
	}

	q := `INSERT INTO push_tokens (user_id, token, platform) 
			VALUES ($1, $2, $3)
			ON CONFLICT (user_id, token) DO UPDATE SET
				platform = EXCLUDED.platform,
				updated_at = NOW()`

	if _, err := tx.Exec(q, userID, token, platform); err != nil {
		return fmt.Errorf("save push token: %w", err)
	}

	return tx.Commit()
}

func GetUserPushTokens(db *sqlx.DB, userID uuid.UUID) ([]PushToken, error) {
	q := `SELECT * FROM push_tokens WHERE user_id = $1`

	var pt []PushToken

	if err := db.Select(&pt, q, userID); err != nil {
		return nil, fmt.Errorf("get push tokens: %w", err)
	}

	return pt, nil
}
