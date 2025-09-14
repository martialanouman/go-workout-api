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

	fmt.Fprintf(w, "This is the workout ID %d\n", workoutId)
}

func (wh WorkoutHandler) HandleCreateWorkout(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Created a workout!\n")
}
