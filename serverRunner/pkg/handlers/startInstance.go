package handlers

import (
	"log"
	"net/http"

	"github.com/3AM-Developer/server-runner/internal/state"
)

func StartHandler(rw http.ResponseWriter, r *http.Request) {
	err := state.AppState.StartInstance()
	if err != nil {
		log.Printf("Error on startHandler: %v", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}
