package app

import (
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joao-vitor-felix/workout-api/internal/api"
	"github.com/joao-vitor-felix/workout-api/internal/store"
	"github.com/joao-vitor-felix/workout-api/migrations"
)

type Application struct {
	Logger         *log.Logger
	WorkoutHandler *api.WorkoutHandler
	DBPool         *pgxpool.Pool
}

func NewApplication() (*Application, error) {
	dbPool, err := store.OpenPool()
	if err != nil {
		return nil, err
	}

	err = store.MigrateFS(dbPool, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	workoutHandler := api.NewWorkoutHandler()
	app := &Application{
		Logger:         logger,
		WorkoutHandler: workoutHandler,
		DBPool:         dbPool,
	}

	return app, nil
}
