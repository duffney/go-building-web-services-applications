package main

import (
	"flag"
	"fmt"
	"net/http"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "dev", "Environment (dev|staging|prod)")
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/healthcheck", healthcheck)

	addr := fmt.Sprintf(":%d", cfg.port)

	err := http.ListenAndServe(addr, mux)
	if err != nil {
		fmt.Println(err)
	}
}

func healthcheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	fmt.Fprintln(w, "status: available")
	fmt.Fprintf(w, "environment: %s\n", "dev") // function cannot access cfg.env
	fmt.Fprintf(w, "version: %s\n", version)
}
