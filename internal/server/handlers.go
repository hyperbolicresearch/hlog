package server

import (
	"fmt"
	"net/http"
)

func LogHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "Hello, log!")
}
