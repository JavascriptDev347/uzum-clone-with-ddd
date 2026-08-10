package http

import "github.com/go-chi/chi/v5"

func NewRouter(h *ProductHandler) chi.Router {
	r := chi.NewRouter()
	// Define routes
	r.Post("/products", h.CreateProduct)

	return r
}
