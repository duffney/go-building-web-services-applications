// rename to handlers.go
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Duffney/go-building-web-services-applications/internal/data"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		app.methodNotAllowedResponse(w, r)
		return
	}
	//fmt.Fprintln(w, "create a new movie")

	var input struct {
		Title   string       `json:"title"`
		Year    int32        `json:"year'`
		Runtime data.Runtime `json:"runtime"`
		Genres  []string     `json:"genres"`
	}

	//err := json.NewDecoder(r.Body).Decode(&input)
	err := app.readJSON(w, r, &input)
	if err != nil {
		//app.errorResponse(w, r, http.StatusBadRequest, err.Error())
		// replace with "specialized helper" for some reason unbeknownest to me
		app.badRequestResponse(w, r, err)
		return
	}
	fmt.Fprintf(w, "%+v\n", input)
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		app.methodNotAllowedResponse(w, r)
		return
	}

	id := r.URL.Path[len("v1/movies//"):]
	idInt, err := strconv.ParseInt(id, 10, 64)

	// TODO add validation for numeric values only

	movie := data.Movie{
		ID:        idInt,
		CreatedAt: time.Now(),
		Title:     "Casablanca",
		Runtime:   102,
		Genres:    []string{"drama", "romance", "war"},
		Version:   1,
	}
	env := envelope{"movie": movie}
	js, err := json.MarshalIndent(env, "", "\t")
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	js = append(js, '\n')
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
