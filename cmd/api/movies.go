// rename to handlers.go
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Duffney/go-building-web-services-applications/internal/data"
	"github.com/Duffney/go-building-web-services-applications/internal/validator"
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

	// validator without helper func
	//v := validator.New()

	//v.Check(input.Title != "", "title", "must be provided")
	//v.Check(len(input.Title) <= 500, "title", "must not be more than 500 bytes long")

	//v.Check(input.Year != 0, "year", "must be provided")
	//v.Check(input.Year >= 1888, "year", "must be grater than 1888")
	//v.Check(input.Year <= int32(time.Now().Year()), "year", "must not be in the future")

	//v.Check(input.Runtime != 0, "runtime", "must be provided")
	//v.Check(input.Runtime > 0, "runtime", "must be a positive integer")

	//v.Check(input.Genres != nil, "genres", "must be provided")
	//v.Check(len(input.Genres) >= 1, "genres", "must contain at least 1 genre")
	//v.Check(len(input.Genres) <= 5, "genres", "must not contain more than 5 genres")
	//v.Check(validator.Unique(input.Genres), "genres", "must not contain duplicate values")

	//if !v.Valid() {
	//	app.failedValidationResponse(w, r, v.Errors)
	//	return
	//}

	movie := &data.Movie{
		Title:   input.Title,
		Year:    input.Year,
		Runtime: input.Runtime,
		Genres:  input.Genres,
	}

	v := validator.New()

	if data.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
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
