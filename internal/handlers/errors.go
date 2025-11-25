package handlers

import (
	"log"
	"net/http"
)

func handleInternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	// TODO: replace with proper logging
	log.Printf("[SERVER ERROR] error: %s, method: %s, path: %s", err.Error(), r.Method, r.URL.Path)

	writeJSONError(w, http.StatusInternalServerError, "internal server error")
}
