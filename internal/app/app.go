package app

import (
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
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

	db := stdlib.OpenDBFromPool(dbPool)
	defer db.Close()
	err = store.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	workoutStore := store.NewPostgresWorkoutStore(stdlib.OpenDBFromPool(dbPool))
	workoutHandler := api.NewWorkoutHandler(workoutStore, logger)
	app := &Application{
		Logger:         logger,
		WorkoutHandler: workoutHandler,
		DBPool:         dbPool,
	}

	return app, nil
}
