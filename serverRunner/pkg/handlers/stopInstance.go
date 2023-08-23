package handlers

import (
	"log"
	"net/http"

	"github.com/3AM-Developer/server-runner/internal/state"
)

func StopHandler(rw http.ResponseWriter, r *http.Request) {
	err := state.AppState.StopInstance()
	if err != nil {
		log.Printf("Error on stopHandler: %v", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}
