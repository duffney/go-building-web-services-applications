package main

import (
	"encoding/json"
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	// v1
	//fmt.Fprintln(w, "status: available")
	//fmt.Fprintf(w, "environment: %s\n", app.config.env)
	//fmt.Fprintf(w, "version: %s\n", version)

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
		app.serverErrorResponse(w, r, err)
		return
	}

	js = append(js, '\n')

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
