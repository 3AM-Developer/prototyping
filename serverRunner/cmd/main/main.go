package main

import (
	"net/http"

	"github.com/3AM-Developer/server-runner/pkg/handlers"
	"github.com/gorilla/mux"
)

func init() {
	// instantiate db variable

	// then I want to pass it to handlers using middleware
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/start", handlers.StartHandler)
	r.HandleFunc("/stop", handlers.StopHandler)
	r.HandleFunc("/update/{id_or_name}", handlers.UpdateHandler)
	// Start the server
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}
