package main

import (
	"fmt"
	"net/http"
)

func (app *application) healthcheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	fmt.Fprintln(w, "status: available")
	fmt.Fprintf(w, "environment: %s\n", app.config.env)
	fmt.Fprintf(w, "version: %s\n", version)
}

func (app *application) getCreateBooksHandler(w http.ResponseWriter, r *http.Request) {
	// if r.Method == http.MethodGet {
	// 	app.getBooks(w, r)
	// 	return
	// }
	// if r.Method == http.MethodPost {
	// 	app.createBook(w, r)
	// 	return
	// }
}

func (app *application) getUpdateDeleteBooksHandler(w http.ResponseWriter, r *http.Request) {
	// switch r.Method {
	// case http.MethodGet:
	// 	app.getBook(w, r)
	// case http.MethodPut:
	// 	app.updateBook(w, r)
	// case http.MethodDelete:
	// 	app.deleteBook(w, r)
	// default:
	// 	http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	// }
}
