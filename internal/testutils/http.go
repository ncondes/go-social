package testutils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// MakeJSONRequest creates an HTTP request with JSON body
func MakeJSONRequest(t *testing.T, method, path string, body any) *http.Request {
	t.Helper()

	var bodyReader io.Reader
	if body != nil {
		jsonBytes, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("failed to marshal request body: %v", err)
		}
		bodyReader = bytes.NewReader(jsonBytes)
	}

	req := httptest.NewRequest(method, path, bodyReader)
	req.Header.Set("Content-Type", "application/json")
	return req
}

// ParseJSONResponse decodes JSON response body into v
func ParseJSONResponse(t *testing.T, w *httptest.ResponseRecorder, v any) {
	t.Helper()

	if err := json.NewDecoder(w.Body).Decode(v); err != nil {
		t.Fatalf("failed to decode response: %v. Body: %s", err, w.Body.String())
	}
}
