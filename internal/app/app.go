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
	UserHandler    *api.UserHandler
	TokenHandler   *api.TokenHandler
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
	//TODO: fix db connection for stores
	workoutStore := store.NewPostgresWorkoutStore(stdlib.OpenDBFromPool(dbPool))
	workoutHandler := api.NewWorkoutHandler(workoutStore, logger)
	userStore := store.NewPostgresUserStore(stdlib.OpenDBFromPool(dbPool))
	userHandler := api.NewUserHandler(userStore, logger)
	tokenStore := store.NewPostgresTokenStore(stdlib.OpenDBFromPool(dbPool))
	tokenHandler := api.NewTokenHandler(tokenStore, userStore, logger)
	app := &Application{
		Logger:         logger,
		WorkoutHandler: workoutHandler,
		UserHandler:    userHandler,
		TokenHandler:   tokenHandler,
		DBPool:         dbPool,
	}
	return app, nil
}
