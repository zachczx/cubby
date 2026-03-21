package server

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/zachczx/cubby/api/internal/gym"
	"github.com/zachczx/cubby/api/internal/response"
)

func (s *Service) NewExerciseHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := s.GetUserIDFromContext(r.Context())
	if err != nil {
		response.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var input gym.ExerciseInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.WriteError(r.Context(), w, err)
		return
	}

	ex, err := gym.NewExercise(s.DB, userID, input)
	if err != nil {
		response.WriteError(r.Context(), w, err)
		return
	}

	response.WriteJSONStatus(r.Context(), w, http.StatusCreated, ex)
}

func (s *Service) GetAllExercisesHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := s.GetUserIDFromContext(r.Context())
	if err != nil {
		response.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	exercises, err := gym.GetAllExercises(s.DB, userID)
	if err != nil {
		response.WriteError(r.Context(), w, err)
		return
	}

	response.WriteJSON(r.Context(), w, exercises)
}

func (s *Service) EditExerciseHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := s.GetUserIDFromContext(r.Context())
	if err != nil {
		response.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	exerciseID, err := uuid.Parse(r.PathValue("exerciseID"))
	if err != nil {
		response.WriteError(r.Context(), w, err)
		return
	}

	var input gym.ExerciseInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.WriteError(r.Context(), w, err)
		return
	}

	if err := gym.EditExercise(s.DB, userID, exerciseID, input); err != nil {
		response.WriteError(r.Context(), w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Service) DeleteExerciseHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := s.GetUserIDFromContext(r.Context())
	if err != nil {
		response.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	exerciseID, err := uuid.Parse(r.PathValue("exerciseID"))
	if err != nil {
		response.WriteError(r.Context(), w, err)
		return
	}

	if err := gym.DeleteExercise(s.DB, userID, exerciseID); err != nil {
		response.WriteError(r.Context(), w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Service) NewWorkoutHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := s.GetUserIDFromContext(r.Context())
	if err != nil {
		response.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	workout, err := gym.NewWorkout(s.DB, userID)
	if err != nil {
		response.WriteError(r.Context(), w, err)
		return
	}

	response.WriteJSONStatus(r.Context(), w, http.StatusCreated, workout)
}

func (s *Service) GetAllWorkoutsHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := s.GetUserIDFromContext(r.Context())
	if err != nil {
		response.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	workouts, err := gym.GetAllWorkouts(s.DB, userID)
	if err != nil {
		response.WriteError(r.Context(), w, err)
		return
	}

	response.WriteJSON(r.Context(), w, workouts)
}

func (s *Service) EditWorkoutHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := s.GetUserIDFromContext(r.Context())
	if err != nil {
		response.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	workoutID, err := uuid.Parse(r.PathValue("workoutID"))
	if err != nil {
		response.WriteError(r.Context(), w, err)
		return
	}

	var input gym.WorkoutInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.WriteError(r.Context(), w, err)
		return
	}

	if err := gym.EditWorkout(s.DB, userID, workoutID, input); err != nil {
		response.WriteError(r.Context(), w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Service) DeleteWorkoutHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := s.GetUserIDFromContext(r.Context())
	if err != nil {
		response.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	workoutID, err := uuid.Parse(r.PathValue("workoutID"))
	if err != nil {
		response.WriteError(r.Context(), w, err)
		return
	}

	if err := gym.DeleteWorkout(s.DB, userID, workoutID); err != nil {
		response.WriteError(r.Context(), w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Service) NewSetHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := s.GetUserIDFromContext(r.Context())
	if err != nil {
		response.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	workoutID, err := uuid.Parse(r.PathValue("workoutID"))
	if err != nil {
		response.WriteError(r.Context(), w, err)
		return
	}

	var input gym.SetInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.WriteError(r.Context(), w, err)
		return
	}

	set, err := gym.NewSet(s.DB, userID, workoutID, input)
	if err != nil {
		response.WriteError(r.Context(), w, err)
		return
	}

	response.WriteJSONStatus(r.Context(), w, http.StatusCreated, set)
}

func (s *Service) EditSetHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := s.GetUserIDFromContext(r.Context())
	if err != nil {
		response.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	setID, err := uuid.Parse(r.PathValue("setID"))
	if err != nil {
		response.WriteError(r.Context(), w, err)
		return
	}

	var input gym.SetInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.WriteError(r.Context(), w, err)
		return
	}

	if err := gym.EditSet(s.DB, userID, setID, input); err != nil {
		response.WriteError(r.Context(), w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Service) DeleteSetHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := s.GetUserIDFromContext(r.Context())
	if err != nil {
		response.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	setID, err := uuid.Parse(r.PathValue("setID"))
	if err != nil {
		response.WriteError(r.Context(), w, err)
		return
	}

	if err := gym.DeleteSet(s.DB, userID, setID); err != nil {
		response.WriteError(r.Context(), w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
