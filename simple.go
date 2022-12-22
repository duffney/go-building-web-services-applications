package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
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

type Movie struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitempty"`
	Runtime   Runtime   `json:"runtime,omitempty"`
	Genres    []string  `json:"genres,omitempty"`
	Version   int32     `json:"version"`
}

type envelope map[string]any

type Runtime int32

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	//data := map[string]string{
	//	"status":      "available",
	//	"environment": app.config.env,
	//	"version":     version,
	//}

	if r.Method != http.MethodGet {
		app.methodNotAllowedResponse(w, r)
		return
	}

	env := envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": app.config.env,
			"version":     version,
		},
	}

	js, err := json.MarshalIndent(env, "", "\t")
	if err != nil {
		//app.logger.Print(err)
		//http.Error(w, "The server encountered a problem", http.StatusInternalServerError)
		app.serverErrorResponse(w, r, err)
		return
	}

	js = append(js, '\n')

	w.Header().Set("Content-Type", "application/json")

	w.Write(js)

	//fmt.Fprintln(w, "status: available")
	//fmt.Fprintf(w, "environment: %s\n", app.config.env)
	//fmt.Fprintf(w, "version: %s\n", version)
}

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	// TODO if http.Method != post err
	fmt.Fprintln(w, "create a new movie")
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	// TODO if http.method != GET err app.methodNotAllowedResponse
	id := r.URL.Path[len("v1/movies//"):]
	idInt, err := strconv.ParseInt(id, 10, 64)

	movie := Movie{
		ID:        idInt,
		CreatedAt: time.Now(),
		Title:     "casablanca",
		Runtime:   102,
		Genres:    []string{"drama", "romance", "war"},
		Version:   1,
	}
	resp := envelope{"movie": movie}
	js, err := json.MarshalIndent(resp, "", "\t")
	if err != nil {
		//app.logger.Print(err)
		//http.Error(w, "The server encountered a problem", http.StatusInternalServerError)
		app.serverErrorResponse(w, r, err)
		//return
	}

	js = append(js, '\n')

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

	//fmt.Fprintf(w, "created movie: %s", id)
	// add validation for numeric values
}

func (app *application) logError(r *http.Request, err error) {
	app.logger.Print(err)
}

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := envelope{"error": message}

	js, err := json.MarshalIndent(env, "", "\t")
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(500)
	}

	js = append(js, '\n')
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
}

func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)

	message := "the server encountered a problem and could not process your request"
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	app.errorResponse(w, r, http.StatusNotFound, message)
}

func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	app.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

func (r Runtime) MarshalJSON() ([]byte, error) {
	jsonValue := fmt.Sprintf("%d mins", r)
	quotedJSONValue := strconv.Quote(jsonValue)
	return []byte(quotedJSONValue), nil
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

	http.HandleFunc("/", app.notFoundResponse) // route for custom 404 (notFoundResponse)
	http.HandleFunc("/v1/healthcheck", app.healthcheckHandler)
	http.HandleFunc("/v1/movies", app.createMovieHandler)
	http.HandleFunc("/v1/movies/", app.showMovieHandler)

	logger.Printf("starting %s server on %s", cfg.env, addr)
	err := http.ListenAndServe(addr, nil)
	logger.Fatal(err)
}
