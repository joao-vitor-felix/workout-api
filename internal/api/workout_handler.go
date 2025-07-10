package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/joao-vitor-felix/workout-api/internal/store"
)

type WorkoutHandler struct {
	store store.WorkoutStore
}

func NewWorkoutHandler(store store.WorkoutStore) *WorkoutHandler {
	return &WorkoutHandler{
		store,
	}
}

func (wh *WorkoutHandler) GetById(w http.ResponseWriter, r *http.Request) {
	workoutIdParam := chi.URLParam(r, "id")
	if workoutIdParam == "" {
		http.NotFound(w, r)
		return
	}

	workoutId, err := strconv.ParseInt(workoutIdParam, 10, 64)

	if err != nil {
		http.Error(w, "Invalid workout ID", http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "this is the workout ID: %d\n", workoutId)
}

func (wh *WorkoutHandler) Create(w http.ResponseWriter, r *http.Request) {
	var workout store.Workout
	err := json.NewDecoder(r.Body).Decode(&workout)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	createdWorkout, err := wh.store.Create(&workout)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create workout: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdWorkout)
}
