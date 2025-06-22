package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type WorkoutHandler struct{}

func NewWorkoutHandler() *WorkoutHandler {
	return &WorkoutHandler{}
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
	fmt.Fprintf(w, "created\n")
}
