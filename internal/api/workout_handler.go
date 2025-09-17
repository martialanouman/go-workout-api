package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/martialanouman/femProject/internal/store"
	"github.com/martialanouman/femProject/internal/utils"
)

type WorkoutHandler struct {
	store store.WorkoutStore
}

func NewWorkoutHandler(store store.WorkoutStore) *WorkoutHandler {
	return &WorkoutHandler{
		store: store,
	}
}

func (wh *WorkoutHandler) HandleGetWorkoutById(w http.ResponseWriter, r *http.Request) {
	workoutId, err := utils.ReadIdParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	workout, err := wh.store.GetWorkoutById(workoutId)
	if err == sql.ErrNoRows {
		http.Error(w, "workout not found", http.StatusNotFound)
		return
	}

	if err != nil {
		fmt.Println(err)
		http.Error(w, "failed to fetch workout", http.StatusNotFound)
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"workout": workout})
}

func (wh WorkoutHandler) HandleCreateWorkout(w http.ResponseWriter, r *http.Request) {
	var workout store.Workout
	err := json.NewDecoder(r.Body).Decode(&workout)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "failed to create workout", http.StatusInternalServerError)
		return
	}

	createdWorkout, err := wh.store.CreateWorkout(&workout)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "failed to create workout", http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"workout": createdWorkout})
}

func (wh *WorkoutHandler) HandleDeleteWorkout(w http.ResponseWriter, r *http.Request) {
	workoutId, err := utils.ReadIdParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	err = wh.store.DeletingWorkout(workoutId)
	if err == sql.ErrNoRows {
		http.Error(w, "workout not found", http.StatusNotFound)
		return
	}

	if err != nil {
		fmt.Println("failed deleting workout: ", err)
		http.Error(w, "failed to delete workout", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (wh *WorkoutHandler) HandleUpdateWorkout(w http.ResponseWriter, r *http.Request) {
	workoutId, err := utils.ReadIdParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	existingWorkout, err := wh.store.GetWorkoutById(workoutId)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "failed to fetch workout", http.StatusInternalServerError)
		return
	}

	if existingWorkout == nil {
		fmt.Println(err)
		http.Error(w, "failed to fetch workout", http.StatusInternalServerError)
		return
	}

	var updateWorkoutRequest struct {
		Title           *string `json:"title"`
		Description     *string `json:"description"`
		DurationMinutes *int    `json:"duration_minutes"`
		CaloriesBurned  *int    `json:"calories_burned"`
		Entries         []store.WorkoutEntry
	}

	err = json.NewDecoder(r.Body).Decode(&updateWorkoutRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if updateWorkoutRequest.Title != nil {
		existingWorkout.Title = *updateWorkoutRequest.Title
	}

	if updateWorkoutRequest.Description != nil {
		existingWorkout.Description = *updateWorkoutRequest.Description
	}

	if updateWorkoutRequest.DurationMinutes != nil {
		existingWorkout.DurationMinutes = *updateWorkoutRequest.DurationMinutes
	}

	if updateWorkoutRequest.CaloriesBurned != nil {
		existingWorkout.CaloriesBurned = *updateWorkoutRequest.CaloriesBurned
	}

	if len(updateWorkoutRequest.Entries) > 0 {
		existingWorkout.Entries = updateWorkoutRequest.Entries
	}

	err = wh.store.UpdateWorkout(existingWorkout)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "failed to update workout", http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"workout": existingWorkout})
}
