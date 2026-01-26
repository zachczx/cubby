package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
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

func WriteJSON(w http.ResponseWriter, data any) {
	bytes, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		WriteError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(bytes); err != nil {
		fmt.Println(err)
		WriteError(w, err)
	}
}

func WriteError(w http.ResponseWriter, err error) {
	var errResp ErrorResponse

	var valErr *ValidationError

	if errors.As(err, &valErr) {
		errResp = ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: valErr.Message,
			Field:   valErr.Field,
		}
		writeJSON(w, http.StatusBadRequest, errResp)
		return
	}

	switch {
	case errors.Is(err, sql.ErrNoRows):
		errResp = ErrorResponse{
			Status:  http.StatusNotFound,
			Message: "could not find resource",
		}
		writeJSON(w, http.StatusNotFound, errResp)

		// Need more cases
	}

	slog.Error("internal error", "error", err)
	errResp = ErrorResponse{
		Status:  http.StatusInternalServerError,
		Message: "an error occurred",
	}

	writeJSON(w, http.StatusInternalServerError, errResp)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	bytes, err := json.Marshal(data)
	if err != nil {
		slog.Error("marshal error", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err := w.Write(bytes); err != nil {
		slog.Error("write response error", "error", err)
	}
}
