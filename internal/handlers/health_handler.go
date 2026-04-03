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

// Check godoc
//
//	@Summary		Check the health of the API
//	@Description	Check the health of the API
//	@Tags			health
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	map[string]string
//	@Router			/health [get]
func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status": "ok",
		"env":    h.config.Env,
	}

	respondWithData(w, http.StatusOK, data)
}
