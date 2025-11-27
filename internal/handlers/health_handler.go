package handlers

import (
	"net/http"

	"github.com/ncondes/go/social/internal/config"
)

type HealthHandler struct {
	config *config.Config
}

func NewHealthHandler(config *config.Config) *HealthHandler {
	return &HealthHandler{
		config: config,
	}
}

func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status": "ok",
		"env":    h.config.Env,
	}

	if err := respondWithData(w, http.StatusOK, data); err != nil {
		handleInternalServerError(w, r, err)
		return
	}
}
