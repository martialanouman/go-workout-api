package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/martialanouman/femProject/internal/middleware"
	"github.com/martialanouman/femProject/internal/store"
	"github.com/martialanouman/femProject/internal/utils"
)

type WorkoutHandler struct {
	store  store.WorkoutStore
	logger *log.Logger
}

func NewWorkoutHandler(store store.WorkoutStore, logger *log.Logger) *WorkoutHandler {
	return &WorkoutHandler{
		store:  store,
		logger: logger,
	}
}

func (h *WorkoutHandler) HandleGetWorkoutById(w http.ResponseWriter, r *http.Request) {
	workoutId, err := utils.ReadIdParam(r)
	if err != nil {
		h.logger.Printf("ERROR: readIdParam %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid workout id"})
		return
	}

	workout, err := h.store.GetWorkoutById(workoutId)
	if err == sql.ErrNoRows {
		h.logger.Printf("ERROR: GetWorkoutById %v", err)
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "workout not found"})
		return
	}

	if err != nil {
		h.logger.Printf("ERROR: GetWorkoutById %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"workout": workout})
}

func (h WorkoutHandler) HandleCreateWorkout(w http.ResponseWriter, r *http.Request) {
	var workout store.Workout
	currentUser := middleware.GetUser(r)

	err := json.NewDecoder(r.Body).Decode(&workout)
	if err != nil {
		h.logger.Printf("ERROR: json.Decode %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request data"})
		return
	}

	workout.UserId = currentUser.Id
	createdWorkout, err := h.store.CreateWorkout(&workout)
	if err != nil {
		h.logger.Printf("ERROR: CreateWorkout %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"workout": createdWorkout})
}

func (h *WorkoutHandler) HandleDeleteWorkout(w http.ResponseWriter, r *http.Request) {
	workoutId, err := utils.ReadIdParam(r)
	if err != nil {
		h.logger.Printf("ERROR: ReadIdParam %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid workout id"})
		return
	}

	currentUser := middleware.GetUser(r)
	ownerId, err := h.store.GetWorkoutOwner(workoutId)
	if errors.Is(err, sql.ErrNoRows) {
		h.logger.Printf("ERROR: GetWorkoutOwner %v", err)
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "workout not found"})
		return
	}

	if err != nil {
		h.logger.Printf("ERROR: GetWorkoutOwner %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	if errors.Is(err, sql.ErrNoRows) {
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "workout not found"})
		return
	}

	if ownerId != currentUser.Id {
		h.logger.Printf("ERROR: unauthorized delete attempt by user %d on workout %d owned by user %d", currentUser.Id, workoutId, ownerId)
		utils.WriteJSON(w, http.StatusForbidden, utils.Envelope{"error": "you do not have permission to delete this workout"})
		return
	}

	err = h.store.DeleteWorkout(workoutId)
	if err == sql.ErrNoRows {
		h.logger.Printf("ERROR: DeletingWorkout %v", err)
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "workout not found"})
		return
	}

	if err != nil {
		h.logger.Printf("ERROR: DeletingWorkout %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *WorkoutHandler) HandleUpdateWorkout(w http.ResponseWriter, r *http.Request) {
	workoutId, err := utils.ReadIdParam(r)
	if err != nil {
		h.logger.Printf("ERROR: ReadIdParam %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid workout id"})
		return
	}

	existingWorkout, err := h.store.GetWorkoutById(workoutId)
	if errors.Is(err, sql.ErrNoRows) {
		h.logger.Printf("ERROR: GetWorkoutById %v", err)
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "workout not found"})
		return
	}

	if err != nil {
		h.logger.Printf("ERROR: GetWorkoutById %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	if existingWorkout == nil {
		h.logger.Printf("ERROR: GetWorkoutById %v", err)
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "workout not found"})
		return
	}

	currentUser := middleware.GetUser(r)
	if existingWorkout.UserId != currentUser.Id {
		h.logger.Printf("ERROR: unauthorized update attempt by user %d on workout %d owned by user %d", currentUser.Id, existingWorkout.Id, existingWorkout.UserId)
		utils.WriteJSON(w, http.StatusForbidden, utils.Envelope{"error": "you do not have permission to update this workout"})
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
		h.logger.Printf("ERROR: validating update payload %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid workout data"})
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

	err = h.store.UpdateWorkout(existingWorkout)
	if err != nil {
		h.logger.Printf("ERROR: UpdateWorkout %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"workout": existingWorkout})
}

func (h *WorkoutHandler) HandleGetWorkouts(w http.ResponseWriter, r *http.Request) {
	take, skip, err := utils.ReadPaginationParams(r)
	if err != nil {
		h.logger.Printf("ERROR: ReadPaginationParams %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid pagination parameters"})
		return
	}

	workouts, err := h.store.GetWorkouts(take, skip)
	if err != nil {
		h.logger.Printf("ERROR: GetWorkouts %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"workouts": workouts, "take": take, "skip": skip})
}
