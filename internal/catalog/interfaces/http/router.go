package http

import (
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/identity/domain"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/identity/infrastructure/security"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/identity/interfaces/http/middleware"
	"github.com/go-chi/chi/v5"
)

func NewRouter(h *ProductHandler, c *CategoryHandler, e *EventHandler, tokenService *security.JWTTokenService) chi.Router {
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
	r.With(middleware.Authenticate(tokenService), middleware.RequireRole(domain.RoleAdmin)).
		Get("/categories/admin", c.GetAllCategories)

	r.Get("/categories", c.GetCategories)
	r.Get("/categories/{id}", c.GetCategory)

	// events (protected)
	r.With(middleware.Authenticate(tokenService), middleware.RequireRole(domain.RoleAdmin)).
		Post("/events", e.CreateEvent)
	r.With(middleware.Authenticate(tokenService), middleware.RequireRole(domain.RoleAdmin)).
		Put("/events/{id}", e.UpdateEvent)
	r.With(middleware.Authenticate(tokenService), middleware.RequireRole(domain.RoleAdmin)).
		Delete("/events/{id}", e.DeleteEvent)
	r.With(middleware.Authenticate(tokenService), middleware.RequireRole(domain.RoleAdmin)).
		Get("/events/admin", e.GetAllEvents)

	r.Get("/events", e.GetEvents)
	r.Get("/events/{id}", e.GetEvent)

	return r
}
