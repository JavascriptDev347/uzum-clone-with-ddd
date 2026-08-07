package http

import (
	"github.com/go-chi/chi/v5"
)

func NewRouter(h *IdentityHandler) chi.Router {
	r := chi.NewRouter()

	// public routes
	r.Post("/register", h.Register)
	r.Post("/login", h.Login)
	r.Post("/refresh", h.Refresh)
	return r
}
