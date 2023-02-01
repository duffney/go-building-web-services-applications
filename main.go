package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func main() {
	r := http.NewServeMux()
	r.HandleFunc("/v1/books/{id}", getBookHandler).Methods("GET")
	http.ListenAndServe(":8080", r)
}

func getBookHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	idInt, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Fprint(w, idInt)
}
