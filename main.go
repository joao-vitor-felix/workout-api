package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/joao-vitor-felix/workout-api/internal/app"
	"github.com/joao-vitor-felix/workout-api/internal/routes"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	var port int
	flag.IntVar(&port, "port", 8080, "Port to run the server on")
	flag.Parse()

	app, err := app.NewApplication()
	if err != nil {
		panic(err)
	}

	defer app.DBPool.Close()
	r := routes.SetupRoutes(app)
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      r,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	app.Logger.Printf("Listening on port %d\n", port)
	err = server.ListenAndServe()

	if err != nil {
		app.Logger.Fatalf("Failed to start server: %v", err)
	}
}
