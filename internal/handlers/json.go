package handlers

import (
	"encoding/json"
	"net/http"
)

func jsonEncode(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func jsonDecode(w http.ResponseWriter, r *http.Request, data any) error {
	const MAX_BYTES = 1_048_576 // 1MB
	// limit the size of the request body to 1MB
	r.Body = http.MaxBytesReader(w, r.Body, MAX_BYTES)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Disallow unknown fields
	return decoder.Decode(data)
}

func respondWithError(w http.ResponseWriter, status int, message string) error {
	type envelope struct {
		Error string `json:"error"`
	}

	return jsonEncode(w, status, &envelope{Error: message})
}

func respondWithErrors(w http.ResponseWriter, status int, errors []string) error {
	type envelope struct {
		Errors []string `json:"errors"`
	}

	return jsonEncode(w, status, &envelope{Errors: errors})
}

func respondWithData(w http.ResponseWriter, status int, data any) error {
	type envelope struct {
		Data any `json:"data"`
	}

	return jsonEncode(w, status, &envelope{Data: data})
}

func respondWithPagination(w http.ResponseWriter, status int, data any, meta any) error {
	type envelope struct {
		Data       any `json:"data"`
		Pagination any `json:"pagination"`
	}

	return jsonEncode(w, status, &envelope{Data: data, Pagination: meta})
}
