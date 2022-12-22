package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
}

type application struct {
	config config
	logger *log.Logger
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	app := &application{
		config: cfg,
		logger: logger,
	}

	addr := fmt.Sprintf(":%d", cfg.port)

	http.HandleFunc("/", app.notFoundResponse) // route for custom 404 response
	http.HandleFunc("/v1/healthcheck", app.healthcheckHandler)
	http.HandleFunc("/v1/movies", app.createMovieHandler)
	http.HandleFunc("/v1/movies/", app.showMovieHandler)

	logger.Printf("starting %s server on %s", cfg.env, addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		logger.Fatal(err)
	}
}
