package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/joao-vitor-felix/workout-api/internal/app"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Route("/workouts", func(r chi.Router) {
		r.Get("/{id}", app.WorkoutHandler.GetById)
		r.Post("/", app.WorkoutHandler.Create)
		r.Put("/{id}", app.WorkoutHandler.UpdateById)
	})
	return r
}
