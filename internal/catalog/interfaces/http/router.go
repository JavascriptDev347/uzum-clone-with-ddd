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

	// categories (protected)
	r.With(middleware.Authenticate(tokenService), middleware.RequireRole(domain.RoleAdmin)).
		Post("/categories", c.CreateCategory)
	r.With(middleware.Authenticate(tokenService), middleware.RequireRole(domain.RoleAdmin)).
		Put("/categories/{id}", c.UpdateCategory)
	r.With(middleware.Authenticate(tokenService), middleware.RequireRole(domain.RoleAdmin)).
		Delete("/categories/{id}", c.DeleteCategory)

	// get categories
	r.Get("/categories", c.GetCategories)
	r.Get("/categories/{id}", c.GetCategory)

	return r
}
