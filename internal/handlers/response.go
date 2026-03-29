package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/ncondes/go/social/internal/dtos"
)

func jsonEncode(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		handleInternalServerError(w, nil, err)
	}
}

func jsonDecode(w http.ResponseWriter, r *http.Request, data any) error {
	const MAX_BYTES = 1_048_576 // 1MB
	// limit the size of the request body to 1MB
	r.Body = http.MaxBytesReader(w, r.Body, MAX_BYTES)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Disallow unknown fields
	err := decoder.Decode(data)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return errors.New("request body is required")
		}

		return err
	}

	return nil
}

func respondWithError(w http.ResponseWriter, status int, message string) {
	type envelope struct {
		Error string `json:"error"`
	}

	jsonEncode(w, status, &envelope{Error: message})
}

func respondWithErrors(w http.ResponseWriter, status int, errors []string) {
	type envelope struct {
		Errors []string `json:"errors"`
	}

	jsonEncode(w, status, &envelope{Errors: errors})
}

func respondWithData(w http.ResponseWriter, status int, data any) {
	type envelope struct {
		Data any `json:"data"`
	}

	jsonEncode(w, status, &envelope{Data: data})
}

func respondWithPaginatedData[T any](w http.ResponseWriter, status int, data []T, pagination dtos.CursorBasedPaginationMeta) {
	type envelope struct {
		Data       []T                            `json:"data"`
		Pagination dtos.CursorBasedPaginationMeta `json:"pagination"`
	}

	jsonEncode(w, status, &envelope{Data: data, Pagination: pagination})
}
