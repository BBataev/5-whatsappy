package endpoints

import "github.com/BBataev/whatsappy/internal/config"

type Handler struct {
	cfg *config.Config
}

func NewHandler(cfg *config.Config) *Handler {
	return &Handler{
		cfg: cfg,
	}
}
