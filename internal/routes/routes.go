package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/martialanouman/femProject/internal/app"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(app.AuthMiddleware.Authenticate)

		r.Get("/workouts/{id}", app.AuthMiddleware.RequireUser(app.WorkoutHandler.HandleGetWorkoutById))
		r.Post("/workouts", app.AuthMiddleware.RequireUser(app.WorkoutHandler.HandleCreateWorkout))
		r.Put("/workouts/{id}", app.AuthMiddleware.RequireUser(app.WorkoutHandler.HandleUpdateWorkout))
		r.Delete("/workouts/{id}", app.AuthMiddleware.RequireUser(app.WorkoutHandler.HandleDeleteWorkout))
		r.Get("/workouts", app.AuthMiddleware.RequireUser(app.WorkoutHandler.HandleGetWorkouts))
	})

	r.Get("/health", app.HealthCheck)

	r.Post("/users", app.UserHandler.HandleRegisterUser)
	r.Post("/tokens/auth", app.TokenHandler.HandleCreateToken)

	return r
}
