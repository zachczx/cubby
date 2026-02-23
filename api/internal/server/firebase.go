package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"slices"
	"time"

	_ "embed"

	"cloud.google.com/go/auth/credentials"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/zachczx/cubby/api/internal/response"
	"google.golang.org/api/option"
)

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

func (f *FCMClient) SendToDevice(ctx context.Context, token, title, body string) (string, error) {
	message := &messaging.Message{
		Token: token,
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
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
	}

	msgID, err := f.client.Send(ctx, message)
	if err != nil {
		return "", err
	}

	return msgID, nil
}

type PushToken struct {
	ID        string    `db:"id"`
	UserID    string    `db:"user_id"`
	Token     string    `db:"token"`
	Platform  *string   `db:"platform"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type PushTokenInput struct {
	Token    string `json:"token"`
	Platform string `json:"platform"`
}

var platforms = []string{"web", "ios", "android"}

func (s *Service) PushTokenHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := s.GetUserIDFromContext(r.Context())
	if err != nil {
		response.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var t PushTokenInput

	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		response.WriteError(w, err)
		return
	}

	if !slices.Contains(platforms, t.Platform) {
		response.RespondWithError(w, http.StatusUnprocessableEntity, "invalid platform")
		return
	}

	if err := SavePushToken(s.DB, userID, t.Token, t.Platform); err != nil {
		response.WriteError(w, err)
		return
	}
}

func SavePushToken(db *sqlx.DB, userID uuid.UUID, token string, platform string) error {
	q := `INSERT INTO push_tokens (user_id, token, platform) 
			VALUES ($1, $2, $3)
			ON CONFLICT (user_id, token) DO UPDATE SET
				platform = EXCLUDED.platform,
				updated_at = NOW()`

	if _, err := db.Exec(q, userID, token, platform); err != nil {
		return fmt.Errorf("save push token: %w", err)
	}

	return nil
}

func (s *Service) GetPushTokens(userID uuid.UUID) ([]PushToken, error) {
	q := `SELECT * FROM push_tokens WHERE user_id = $1`

	var pt []PushToken

	if err := s.DB.Select(&pt, q, userID); err != nil {
		return nil, fmt.Errorf("get push tokens: %w", err)
	}

	return pt, nil
}

func (s *Service) NotificationHandler(w http.ResponseWriter, r *http.Request) {
	fcmClient, err := NewFCMClient(r.Context())
	if err != nil {
		response.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	userID, err := s.GetUserIDFromContext(r.Context())
	if err != nil {
		fmt.Println("here")
		response.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	tokens, err := s.GetPushTokens(userID)
	if err != nil {
		response.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if _, err := fcmClient.SendToDevice(r.Context(),
		tokens[0].Token,
		"Testing Title", "test body"); err != nil {
		response.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
}
