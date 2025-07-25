package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/joao-vitor-felix/workout-api/internal/middleware"
	"github.com/joao-vitor-felix/workout-api/internal/store"
	"github.com/joao-vitor-felix/workout-api/internal/utils"
)

type WorkoutHandler struct {
	store  store.WorkoutStore
	logger *log.Logger
}

func NewWorkoutHandler(store store.WorkoutStore, logger *log.Logger) *WorkoutHandler {
	return &WorkoutHandler{
		store,
		logger,
	}
}

func (wh *WorkoutHandler) GetById(w http.ResponseWriter, r *http.Request) {
	workoutId, err := utils.ReadIdParam(r)
	if err != nil {
		wh.logger.Printf("ERROR: reading workout ID: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{
			"error": "Invalid workout ID",
		})
		return
	}

	workout, err := wh.store.GetByID(workoutId)
	if err != nil {
		wh.logger.Printf("ERROR: get workout by ID: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{
			"error": "internal server error",
		})
		return
	}

	if workout == nil {
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{
			"error": "not found",
		})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{
		"data": workout,
	})
}

func (wh *WorkoutHandler) Create(w http.ResponseWriter, r *http.Request) {
	var workout store.Workout
	err := json.NewDecoder(r.Body).Decode(&workout)
	if err != nil {
		wh.logger.Printf("ERROR: invalid body: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{
			"error": "invalid request body",
		})
		return
	}

	currentUser := middleware.GetUser(r)
	if currentUser == nil || currentUser == store.AnonymousUser {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "you must bed logged in"})
		return
	}

	workout.UserID = currentUser.ID

	createdWorkout, err := wh.store.Create(&workout)
	if err != nil {
		wh.logger.Printf("ERROR: create workout: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{
			"error": "internal server error",
		})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{
		"data": createdWorkout,
	})
}

func (wh *WorkoutHandler) UpdateById(w http.ResponseWriter, r *http.Request) {
	workoutId, err := utils.ReadIdParam(r)
	if err != nil {
		wh.logger.Printf("ERROR: reading workout ID: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{
			"error": "Invalid workout ID",
		})
		return
	}

	workout, err := wh.store.GetByID(workoutId)
	if err != nil {
		wh.logger.Printf("ERROR: get workout: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{
			"error": "internal server error",
		})
		return
	}

	if workout == nil {
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{
			"error": "not found",
		})
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
		wh.logger.Printf("ERROR: decoding: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{
			"error": "invalid request body",
		})
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

	currentUser := middleware.GetUser(r)
	if currentUser == nil || currentUser == store.AnonymousUser {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "you must be logged in to update"})
		return
	}

	workoutOwner, err := wh.store.GetWorkoutOwner(workoutId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "workout does not exist"})
			return
		}

		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	if workoutOwner != currentUser.ID {
		utils.WriteJSON(w, http.StatusForbidden, utils.Envelope{"error": "you are not authorized to update this workout"})
		return
	}

	err = wh.store.Update(workout)
	if err != nil {
		wh.logger.Printf("ERROR: update workout: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{
			"error": "internal server error",
		})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{
		"data": workout,
	})
}

func (wh *WorkoutHandler) DeleteById(w http.ResponseWriter, r *http.Request) {
	workoutId, err := utils.ReadIdParam(r)
	if err != nil {
		wh.logger.Printf("ERROR: reading workout ID: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{
			"error": "Invalid workout ID",
		})
		return
	}

	currentUser := middleware.GetUser(r)
	if currentUser == nil || currentUser == store.AnonymousUser {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "you must be logged in to update"})
		return
	}

	workoutOwner, err := wh.store.GetWorkoutOwner(workoutId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "workout does not exist"})
			return
		}

		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	if workoutOwner != currentUser.ID {
		utils.WriteJSON(w, http.StatusForbidden, utils.Envelope{"error": "you are not authorized to update this workout"})
		return
	}

	err = wh.store.Delete(workoutId)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{
				"error": "not found",
			})
			return
		}
		wh.logger.Printf("ERROR: delete workout: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{
			"error": "internal server error",
		})
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, nil)
}
