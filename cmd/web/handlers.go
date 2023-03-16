package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	books, err := app.readinglist.GetAll()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "<html><head><title>Reading List</title></head><body><h1>Reading List</h1><ul>")
	for _, book := range *books {
		fmt.Fprintf(w, "<li>%s (%d)</li>", book.Title, book.Pages)
	}
	fmt.Fprintf(w, "</ul></body></html>")
}

func (app *application) bookView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	book, err := app.readinglist.Get(int64(id))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%s (%d)\n", book.Title, book.Pages)
}

func (app *application) bookCreate(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		app.bookCreateForm(w, r)
	case http.MethodPost:
		app.bookCreateProcess(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}

func (app *application) bookCreateForm(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<html><head><title>Create Book</title></head>"+
		"<body><h1>Create Book</h1><form action=\"/book/create\" method=\"post\">"+
		"<label for=\"title\">Title</label><input type=\"text\" name=\"title\" id=\"title\">"+
		"<label for=\"pages\">Pages</label><input type=\"number\" name=\"pages\" id=\"pages\">"+
		"<label for=\"published\">Published</label><input type=\"number\" name=\"published\" id=\"published\">"+
		"<label for=\"genres\">Genres</label><input type=\"text\" name=\"genres\" id=\"genres\">"+
		"<label for=\"rating\">Rating</label><input type=\"number\" step=\"0.1\" name=\"rating\" id=\"rating\">"+
		"<button type=\"submit\">Create</button></form></body></html>")
}

func (app *application) bookCreateProcess(w http.ResponseWriter, r *http.Request) {
	title := r.PostFormValue("title")
	if title == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	pages, err := strconv.Atoi(r.PostFormValue("pages"))
	if err != nil || pages < 1 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	published, err := strconv.Atoi(r.PostFormValue("published"))
	if err != nil || published < 1 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	genres := strings.Split(r.PostFormValue("genres"), " ")

	ratingFloat, err := strconv.ParseFloat(r.PostFormValue("rating"), 32)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	rating := float32(ratingFloat)

	book := struct {
		Title     string   `json:"title"`
		Pages     int      `json:"pages"`
		Published int      `json:"published"`
		Genres    []string `json:"genres"`
		Rating    float32  `json:"rating"`
	}{
		Title:     title,
		Pages:     pages,
		Published: published,
		Genres:    genres,
		Rating:    rating,
	}

	data, err := json.Marshal(book)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	req, _ := http.NewRequest("POST", app.readinglist.Endpoint, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		log.Printf("unexpected status: %s", resp.Status)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
