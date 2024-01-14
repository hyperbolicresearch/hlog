package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hyperbolicresearch/hlog/internal/core"
)

func LogHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		payload := &core.LogPayloadFromRequest{};
		if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Err occured when reading the req body: %v", err)
		}
		
		// 1. validate the payload
		// 2. store to the db

		// Just write the header if everything goes the right way.
		w.WriteHeader(http.StatusOK)
	}
}
