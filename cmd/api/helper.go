package main

import (
	"encoding/json"
	"net/http"
)

type envelope map[string]any

// Credit: Alex Edwards, Let's Go Further
func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	js = append(js, '\n')

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}
