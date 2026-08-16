package http

import (
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/identity/infrastructure/security"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/identity/interfaces/http/middleware"
	"github.com/go-chi/chi/v5"
)

func NewRouter(h *IdentityHandler, tokenService *security.JWTTokenService) chi.Router {
	r := chi.NewRouter()

	// public routes
	r.Post("/register", h.Register)
	r.Post("/login", h.Login)
	r.Post("/refresh", h.Refresh)

	// protected routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.Authenticate(tokenService))
		r.Get("/me", h.GetMe)
	})

	return r
}
