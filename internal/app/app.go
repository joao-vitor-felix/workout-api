package app

import (
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joao-vitor-felix/workout-api/internal/api"
	"github.com/joao-vitor-felix/workout-api/internal/store"
)

type Application struct {
	Logger         *log.Logger
	WorkoutHandler *api.WorkoutHandler
	DBPool         *pgxpool.Pool
}

func NewApplication() (*Application, error) {
	db, err := store.Open()

	if err != nil {
		return nil, err
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	workoutHandler := api.NewWorkoutHandler()

	app := &Application{
		Logger:         logger,
		WorkoutHandler: workoutHandler,
		DBPool:         db,
	}

	return app, nil
}
