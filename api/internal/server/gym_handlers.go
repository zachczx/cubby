package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/zachczx/cubby/api/internal/gym"
	"github.com/zachczx/cubby/api/internal/response"
)

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

func (s *Service) ReorderSetHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := s.GetUserIDFromContext(r.Context())
	if err != nil {
		response.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var input gym.ReorderSetInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.WriteError(r.Context(), w, err)
		return
	}

	if err := gym.ReorderSet(s.DB, userID, input); err != nil {
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

func (s *Service) GetFavouritesHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := s.GetUserIDFromContext(r.Context())
	if err != nil {
		response.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	ids, err := gym.GetFavourites(s.DB, userID)
	if err != nil {
		response.WriteError(r.Context(), w, err)
		return
	}

	response.WriteJSON(r.Context(), w, map[string][]string{"exerciseIds": ids})
}

func (s *Service) ToggleFavouriteHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := s.GetUserIDFromContext(r.Context())
	if err != nil {
		response.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var input struct {
		ExerciseID string `json:"exerciseId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.WriteError(r.Context(), w, err)
		return
	}

	ids, err := gym.ToggleFavourite(s.DB, userID, input.ExerciseID)
	if err != nil {
		response.WriteError(r.Context(), w, err)
		return
	}

	response.WriteJSON(r.Context(), w, map[string][]string{"exerciseIds": ids})
}

// Routine handlers

func (s *Service) GetAllRoutinesHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := s.GetUserIDFromContext(r.Context())
	if err != nil {
		response.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	routines, err := gym.GetAllRoutines(s.DB, userID)
	if err != nil {
		response.WriteError(r.Context(), w, err)
		return
	}

	response.WriteJSON(r.Context(), w, routines)
}

func (s *Service) NewRoutineHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := s.GetUserIDFromContext(r.Context())
	if err != nil {
		response.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var input gym.RoutineInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.WriteError(r.Context(), w, err)
		return
	}

	routine, err := gym.NewRoutine(s.DB, userID, input)
	if err != nil {
		response.WriteError(r.Context(), w, err)
		return
	}

	response.WriteJSONStatus(r.Context(), w, http.StatusCreated, routine)
}

func (s *Service) EditRoutineHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := s.GetUserIDFromContext(r.Context())
	if err != nil {
		response.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	routineID, err := uuid.Parse(r.PathValue("routineID"))
	if err != nil {
		response.WriteError(r.Context(), w, err)
		return
	}

	var input gym.RoutineInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.WriteError(r.Context(), w, err)
		return
	}

	if err := gym.EditRoutine(s.DB, userID, routineID, input); err != nil {
		response.WriteError(r.Context(), w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Service) DeleteRoutineHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := s.GetUserIDFromContext(r.Context())
	if err != nil {
		response.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	routineID, err := uuid.Parse(r.PathValue("routineID"))
	if err != nil {
		response.WriteError(r.Context(), w, err)
		return
	}

	if err := gym.DeleteRoutine(s.DB, userID, routineID); err != nil {
		response.WriteError(r.Context(), w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Service) AddRoutineExerciseHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := s.GetUserIDFromContext(r.Context())
	if err != nil {
		response.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	routineID, err := uuid.Parse(r.PathValue("routineID"))
	if err != nil {
		response.WriteError(r.Context(), w, err)
		return
	}

	var input gym.RoutineExerciseInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.WriteError(r.Context(), w, err)
		return
	}

	exercise, err := gym.AddRoutineExercise(s.DB, userID, routineID, input)
	if err != nil {
		response.WriteError(r.Context(), w, err)
		return
	}

	response.WriteJSONStatus(r.Context(), w, http.StatusCreated, exercise)
}

func (s *Service) EditRoutineExerciseHandler(w http.ResponseWriter, r *http.Request) {
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

	var input gym.RoutineExerciseInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.WriteError(r.Context(), w, err)
		return
	}

	if err := gym.EditRoutineExercise(s.DB, userID, exerciseID, input); err != nil {
		response.WriteError(r.Context(), w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Service) RemoveRoutineExerciseHandler(w http.ResponseWriter, r *http.Request) {
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

	if err := gym.RemoveRoutineExercise(s.DB, userID, exerciseID); err != nil {
		response.WriteError(r.Context(), w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Service) ReorderRoutineExerciseHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := s.GetUserIDFromContext(r.Context())
	if err != nil {
		response.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var input gym.ReorderRoutineExerciseInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.WriteError(r.Context(), w, err)
		return
	}

	if err := gym.ReorderRoutineExercise(s.DB, userID, input); err != nil {
		response.WriteError(r.Context(), w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Service) StartWorkoutFromRoutineHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := s.GetUserIDFromContext(r.Context())
	if err != nil {
		response.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	routineID, err := uuid.Parse(r.PathValue("routineID"))
	if err != nil {
		response.WriteError(r.Context(), w, err)
		return
	}

	workout, err := gym.StartWorkoutFromRoutine(s.DB, userID, routineID)
	if err != nil {
		response.WriteError(r.Context(), w, err)
		return
	}

	response.WriteJSONStatus(r.Context(), w, http.StatusCreated, workout)
}

func (s *Service) GetGymSummaryHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := s.GetUserIDFromContext(r.Context())
	if err != nil {
		response.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	summary, err := gym.GetSummary(s.DB, userID)
	if err != nil {
		response.WriteError(r.Context(), w, err)
		return
	}

	response.WriteJSON(r.Context(), w, summary)
}

func (s *Service) GetGymCalendarHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := s.GetUserIDFromContext(r.Context())
	if err != nil {
		response.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	entries, err := gym.GetCalendarWorkouts(s.DB, userID)
	if err != nil {
		response.WriteError(r.Context(), w, err)
		return
	}

	response.WriteJSON(r.Context(), w, entries)
}

func (s *Service) GetGymMusclesHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := s.GetUserIDFromContext(r.Context())
	if err != nil {
		response.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	weeks := 4
	if wVal := r.URL.Query().Get("weeks"); wVal != "" {
		if w, err := strconv.Atoi(wVal); err == nil && w > 0 {
			weeks = w
		}
	}

	stats, err := gym.GetMusclesFailureStats(s.DB, userID, weeks)
	if err != nil {
		response.WriteError(r.Context(), w, err)
		return
	}

	response.WriteJSON(r.Context(), w, stats)
}

func (s *Service) GetGymUserExercisesHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := s.GetUserIDFromContext(r.Context())
	if err != nil {
		response.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	exercises, err := gym.GetUserExercises(s.DB, userID)
	if err != nil {
		response.WriteError(r.Context(), w, err)
		return
	}

	response.WriteJSON(r.Context(), w, exercises)
}

func (s *Service) GetGymExerciseStatsHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := s.GetUserIDFromContext(r.Context())
	if err != nil {
		response.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	exerciseID := r.PathValue("exerciseID")
	if exerciseID == "" {
		response.RespondWithError(w, http.StatusBadRequest, "exercise_id is required")
		return
	}

	stats, err := gym.GetExerciseStats(s.DB, userID, exerciseID)
	if err != nil {
		response.WriteError(r.Context(), w, err)
		return
	}

	response.WriteJSON(r.Context(), w, stats)
}
