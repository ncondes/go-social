package pagination

import (
	"encoding/base64"
	"encoding/json"
)

func EncodeCursor[T any](cursor T) (string, error) {
	// Serialize the cursor into a JSON byte slice
	data, err := json.Marshal(cursor)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(data), nil
}

func DecodeCursor[T any](encoded string) (T, error) {
	var zero T

	if encoded == "" {
		return zero, nil
	}

	data, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		return zero, err
	}

	var cursor T

	if err := json.Unmarshal(data, &cursor); err != nil {
		return zero, err
	}

	return cursor, nil
}
