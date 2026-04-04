package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/ncondes/go/social/internal/dtos"
	"github.com/ncondes/go/social/internal/logging"
)

func jsonEncode(w http.ResponseWriter, status int, data any, logger logging.Logger) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		handleInternalServerError(w, nil, err, logger)
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

func respondWithError(w http.ResponseWriter, status int, message string, logger logging.Logger) {
	jsonEncode(w, status, &dtos.ErrorResponseDTO{Error: message}, logger)
}

func respondWithErrors(w http.ResponseWriter, status int, errors []string, logger logging.Logger) {
	jsonEncode(w, status, &dtos.ErrorsResponseDTO{Errors: errors}, logger)
}

func respondWithData(w http.ResponseWriter, status int, data any, logger logging.Logger) {
	jsonEncode(w, status, &dtos.DataResponseDTO{Data: data}, logger)
}

func respondWithPaginatedData[T any](w http.ResponseWriter, status int, data []T, pagination dtos.CursorBasedPaginationMetaDTO, logger logging.Logger) {
	jsonEncode(w, status, &dtos.CursorBasedPaginationResponseDTO[T]{Data: data, Pagination: pagination}, logger)
}
