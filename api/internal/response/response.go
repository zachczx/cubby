package response

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/zachczx/cubby/api/internal/logging"
)

type IDResponse struct {
	ID uuid.UUID `json:"id"`
}

type ErrorResponse struct {
	Status  int               `json:"error,omitempty"`
	Message string            `json:"message,omitempty"`
	Field   string            `json:"field,omitempty"`
	Errors  map[string]string `json:"errors,omitempty"`
}

func WriteJSON(ctx context.Context, w http.ResponseWriter, data any) {
	bytes, err := json.Marshal(data)
	if err != nil {
		logging.Error(ctx, "marshal error", "err", err)
		WriteError(ctx, w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if _, err := w.Write(bytes); err != nil {
		logging.Error(ctx, "write bytes error", "error", err)
		WriteError(ctx, w, err)
	}
}

func WriteJSONStatus(ctx context.Context, w http.ResponseWriter, status int, data any) {
	bytes, err := json.Marshal(data)
	if err != nil {
		logging.Error(ctx, "marshal error", "err", err)
		WriteError(ctx, w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if _, err := w.Write(bytes); err != nil {
		logging.Error(ctx, "write bytes error", "error", err)
		WriteError(ctx, w, err)
	}
}

func WriteError(ctx context.Context, w http.ResponseWriter, err error) {
	var errResp ErrorResponse

	var valErr *ValidationError

	if errors.As(err, &valErr) {
		errResp = ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: valErr.Message,
			Field:   valErr.Field,
		}
		writeJSON(ctx, w, http.StatusBadRequest, errResp)
		return
	}

	var pgErr *pgconn.PgError

	switch {
	case errors.Is(err, sql.ErrNoRows):
		errResp = ErrorResponse{
			Status:  http.StatusNotFound,
			Message: "could not find resource",
		}
		writeJSON(ctx, w, http.StatusNotFound, errResp)
		return

	case errors.As(err, &pgErr) && pgErr.Code == "23505":
		errResp = ErrorResponse{
			Status:  http.StatusConflict,
			Message: "record already exists",
		}
		writeJSON(ctx, w, http.StatusConflict, errResp)
		return
	}

	logging.Error(ctx, "internal error", "error", err)
	errResp = ErrorResponse{
		Status:  http.StatusInternalServerError,
		Message: "an error occurred",
	}

	writeJSON(ctx, w, http.StatusInternalServerError, errResp)
}

func writeJSON(ctx context.Context, w http.ResponseWriter, status int, data any) {
	bytes, err := json.Marshal(data)
	if err != nil {
		logging.Error(ctx, "marshal error", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err := w.Write(bytes); err != nil {
		logging.Error(ctx, "write response error", "error", err)
	}
}

func RespondWithError(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	fmt.Fprintln(w, message)
}
