package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/hyperbolicresearch/hlog/internal/core"
)

var validate *validator.Validate

func LogHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		payload := &core.LogPayloadFromRequest{};
		if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Err occured when reading the req body: %v", err)
		}

		validate = validator.New(validator.WithRequiredStructEnabled())
		err := validate.Struct(payload)
		if err != nil {
			fmt.Fprintf(w, "Invalid payload: %v", err)
		}

		// 2. store to the db

		// Just write the header if everything goes the right way.
		w.WriteHeader(http.StatusOK)
	}
}
