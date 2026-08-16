package http

import (
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/identity/domain"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/identity/infrastructure/security"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/identity/interfaces/http/middleware"
	"github.com/go-chi/chi/v5"
)

func NewRouter(h *ProductHandler, c *CategoryHandler, tokenService *security.JWTTokenService) chi.Router {
	r := chi.NewRouter()
	// Define routes
	r.With(middleware.Authenticate(tokenService), middleware.RequireRole(domain.RoleAdmin)).
		Post("/products", h.CreateProduct)

	// categories
	r.Post("/categories", c.CreateCategory)

	return r
}
