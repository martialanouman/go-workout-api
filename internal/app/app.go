package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/martialanouman/femProject/internal/api"
	"github.com/martialanouman/femProject/internal/middleware"
	"github.com/martialanouman/femProject/internal/store"
	"github.com/martialanouman/femProject/migrations"
)

type Application struct {
	Logger         *log.Logger
	WorkoutHandler *api.WorkoutHandler
	UserHandler    *api.UserHandler
	TokenHandler   *api.TokenHandler
	AuthMiddleware middleware.UserMiddleware
	Db             *sql.DB
}

func NewApplication() (*Application, error) {
	db, err := store.Open()
	if err != nil {
		return nil, err
	}

	err = store.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	userStore := store.NewPostgresUserStore(db)

	app := &Application{
		Logger:         logger,
		WorkoutHandler: api.NewWorkoutHandler(store.NewPostgresWorkoutStore(db), logger),
		UserHandler:    api.NewUserHandler(userStore, logger),
		TokenHandler:   api.NewTokenHandler(store.NewPostgresTokenStore(db), userStore, logger),
		AuthMiddleware: middleware.UserMiddleware{Store: userStore},
		Db:             db,
	}

	return app, nil
}

func (a *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Status is available.\n")
}
