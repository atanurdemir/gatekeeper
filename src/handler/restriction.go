package handler

import (
	"encoding/json"
	"net/http"

	"github.com/atanurdemir/gatekeeper/src/service"
)

func RestrictionHandler(w http.ResponseWriter, r *http.Request) {
	state, err := service.Waf.GetDeniedState(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(state)
}
