package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/martialanouman/femProject/internal/store"
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
	paramWorkoutId := chi.URLParam(r, "id")
	if paramWorkoutId == "" {
		http.NotFound(w, r)
		return
	}

	workoutId, err := strconv.ParseInt(paramWorkoutId, 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	workout, err := wh.store.GetWorkoutById(workoutId)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "failed to fetch workout", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-TYpe", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(workout)
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdWorkout)
}

func (wh *WorkoutHandler) HandleUpdateWorkout(w http.ResponseWriter, r *http.Request) {
	paramWorkoutId := chi.URLParam(r, "id")
	if paramWorkoutId == "" {
		http.NotFound(w, r)
		return
	}

	workoutId, err := strconv.ParseInt(paramWorkoutId, 10, 64)
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(existingWorkout)
}
