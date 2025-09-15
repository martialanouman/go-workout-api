package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/martialanouman/femProject/internal/api"
	"github.com/martialanouman/femProject/internal/store"
	"github.com/martialanouman/femProject/migrations"
)

type Application struct {
	Logger         *log.Logger
	WorkoutHandler *api.WorkoutHandler
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

	app := &Application{
		Logger:         logger,
		WorkoutHandler: api.NewWorkoutHandler(store.NewPostgresWorkoutStore(db)),
		Db:             db,
	}

	return app, nil
}

func (a *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Status is available.\n")
}
