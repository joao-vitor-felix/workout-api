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

	workout, err := wh.store.GetByID(workoutId)

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get workout: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(workout)
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

func (wh *WorkoutHandler) UpdateById(w http.ResponseWriter, r *http.Request) {
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

	workout, err := wh.store.GetByID(workoutId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get workout: %v", err), http.StatusInternalServerError)
		return
	}

	if workout == nil {
		http.NotFound(w, r)
		return
	}

	var updateWorkout struct {
		Title           *string              `json:"title"`
		Description     *string              `json:"description"`
		DurationMinutes *int                 `json:"duration_minutes"`
		CaloriesBurned  *int                 `json:"calories_burned"`
		Entries         []store.WorkoutEntry `json:"entries"`
	}

	err = json.NewDecoder(r.Body).Decode(&updateWorkout)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if updateWorkout.Title != nil {
		workout.Title = *updateWorkout.Title
	}
	if updateWorkout.Description != nil {
		workout.Description = *updateWorkout.Description
	}
	if updateWorkout.DurationMinutes != nil {
		workout.DurationMinutes = *updateWorkout.DurationMinutes
	}
	if updateWorkout.CaloriesBurned != nil {
		workout.CaloriesBurned = *updateWorkout.CaloriesBurned
	}
	if updateWorkout.Entries != nil {
		workout.Entries = updateWorkout.Entries
	}

	err = wh.store.Update(workout)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update workout: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(workout)
}
