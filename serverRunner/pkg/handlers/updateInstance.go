package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"path"

	"github.com/3AM-Developer/server-runner/internal/instance"
	"github.com/3AM-Developer/server-runner/internal/models"
	"github.com/3AM-Developer/server-runner/internal/state"
)

var (
	ErrorInstanceNotDefined error
)

type payloadRequest struct {
	Name string `json:"name,omitempty"`
	Id   *int   `json:"id,omitempty"`
}

func UpdateHandler(rw http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	plr := &payloadRequest{}
	if err := decoder.Decode(plr); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	// We need to handle the options of either name or Id.

	var inst *instance.Instance
	var err error
	if plr.Id != nil {
		inst, err = models.InstanceDB.GetInstanceById(*plr.Id)
	}

	if plr.Name != "" {
		inst, err = models.InstanceDB.GetInstanceByName(plr.Name)
	}

	if err != nil || inst == nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	instance.LoadJson(inst, path.Join(inst.Dir, "instance.json"))

	if ok, err := state.AppState.RegisterInstance(inst); !ok {
		if err == state.ErrorInstnaceInvalid {
			log.Printf("Error on start handler: %v", state.ErrorInstnaceInvalid)
			http.Error(rw, state.ErrorInstnaceInvalid.Error(), http.StatusBadRequest)
			return
		}
		log.Printf("Error on startHandler: %v", ErrorInstanceNotDefined)
		http.Error(rw, ErrorInstanceNotDefined.Error(), http.StatusInternalServerError)
		return
	}
}
