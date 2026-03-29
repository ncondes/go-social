package pagination

import (
	"testing"
	"time"
)

func TestEncodeCursor(t *testing.T) {
	t.Run("should encode cursor string", func(t *testing.T) {
		cursor := "cursor"
		encoded, err := EncodeCursor(cursor)
		if err != nil {
			t.Fatal(err)
		}

		expected := "ImN1cnNvciI="

		if encoded != expected {
			t.Errorf("expected %s, got %s", expected, encoded)
		}
	})

	t.Run("should encode a cursor struct", func(t *testing.T) {
		type Cursor struct {
			ID        int64
			CreatedAt time.Time
		}

		cursor := Cursor{
			ID:        1,
			CreatedAt: time.Date(2025, 10, 10, 10, 10, 10, 10, time.UTC),
		}

		encoded, err := EncodeCursor(cursor)
		if err != nil {
			t.Fatal(err)
		}

		expected := "eyJJRCI6MSwiQ3JlYXRlZEF0IjoiMjAyNS0xMC0xMFQxMDoxMDoxMC4wMDAwMDAwMVoifQ=="

		if encoded != expected {
			t.Errorf("expected %s, got %s", expected, encoded)
		}
	})

	t.Run("should return error if cursor is not JSON serializable", func(t *testing.T) {
		cursor := make(chan int) // Channels are not JSON serializable
		_, err := EncodeCursor(cursor)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestDecodeCursor(t *testing.T) {
	t.Run("should decode encoded cursor string", func(t *testing.T) {
		encoded := "ImN1cnNvciI="
		cursor, err := DecodeCursor[string](encoded)
		if err != nil {
			t.Fatal(err)
		}

		expected := "cursor"

		if cursor != expected {
			t.Errorf("expected %s, got %s", expected, cursor)
		}
	})

	t.Run("should decode encoded cursor struct", func(t *testing.T) {
		type Cursor struct {
			ID        int64
			CreatedAt time.Time
		}

		encoded := "eyJJRCI6MSwiQ3JlYXRlZEF0IjoiMjAyNS0xMC0xMFQxMDoxMDoxMC4wMDAwMDAwMVoifQ=="
		cursor, err := DecodeCursor[Cursor](encoded)
		if err != nil {
			t.Fatal(err)
		}

		expected := Cursor{
			ID:        1,
			CreatedAt: time.Date(2025, 10, 10, 10, 10, 10, 10, time.UTC),
		}

		if cursor != expected {
			t.Errorf("expected %+v, got %+v", expected, cursor)
		}
	})

	t.Run("should return zero value if encoded cursor is empty", func(t *testing.T) {
		cursor, err := DecodeCursor[string]("")
		if err != nil {
			t.Fatal(err)
		}

		if cursor != "" {
			t.Errorf("expected empty string, got %s", cursor)
		}
	})

	t.Run("should return error if encoded cursor is not valid", func(t *testing.T) {
		encoded := "something that is clearly not encoded..."
		_, err := DecodeCursor[string](encoded)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("should return error id decoded data is not valid JSON", func(t *testing.T) {
		encoded := "ImN1cnNvciI=" // base64 encoded "cursor"

		type Cursor struct {
			ID int64
		}

		_, err := DecodeCursor[Cursor](encoded)
		if err == nil {
			t.Fatal("expected JSON unmarshal error, got nil")
		}
	})
}
