package http

import "github.com/go-chi/chi/v5"

func NewRouter(h *ProductHandler, c *CategoryHandler) chi.Router {
	r := chi.NewRouter()
	// Define routes
	r.Post("/products", h.CreateProduct)

	// categories
	r.Post("/categories", c.CreateCategory)

	return r
}
