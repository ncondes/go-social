package handlers

import "net/http"

type Health struct{}

func (h *Health) Check(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
