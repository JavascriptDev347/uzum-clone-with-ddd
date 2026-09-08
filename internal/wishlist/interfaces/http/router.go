package http

import (
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/identity/infrastructure/security"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/identity/interfaces/http/middleware"
	"github.com/go-chi/chi/v5"
)

func NewRouter(h *WishlistHandler, tokenService *security.JWTTokenService) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.Authenticate(tokenService))

	r.Get("/", h.GetWishlist)
	r.Post("/items/{product_id}", h.AddItem)
	r.Delete("/items/{product_id}", h.RemoveItem)

	return r
}
