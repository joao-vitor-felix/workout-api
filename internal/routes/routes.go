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
		r.Delete("/{id}", app.WorkoutHandler.DeleteById)
	})
	r.Route("/users", func(r chi.Router) {
		r.Post("/", app.UserHandler.RegisterUser)
	})
	r.Route("/auth", func(r chi.Router) {
		r.Post("/sign-in", app.TokenHandler.Create)
	})
	return r
}
