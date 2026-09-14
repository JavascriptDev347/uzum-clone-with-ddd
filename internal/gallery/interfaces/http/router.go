package http

import (
	identitydomain "github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/identity/domain"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/identity/infrastructure/security"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/identity/interfaces/http/middleware"
	"github.com/go-chi/chi/v5"
)

// NewGalleryRouter - GET /api/v1/gallery uchun router. Ochiq - har kim galereyani ko'rishi
// mumkin, autentifikatsiya shart emas.
func NewGalleryRouter(h *GalleryHandler) chi.Router {
	r := chi.NewRouter()

	r.Get("/", h.ListGalleryPosts)

	return r
}

// NewAdminGalleryRouter - POST /api/v1/admin/gallery va DELETE /api/v1/admin/gallery/{id}
// uchun router. "/api/v1/admin/gallery" va "/api/v1/gallery" ikkalasi ham bare "/api/v1"dan
// farqli, aniq prefikslar bo'lgani uchun (ordering/review'da bo'lgani kabi) alohida-alohida
// mount qilinishi mumkin - Catalog'ning bare "/api/v1"ni egallashi bilan to'qnashmaydi.
func NewAdminGalleryRouter(h *GalleryHandler, tokenService *security.JWTTokenService) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.Authenticate(tokenService))
	r.Use(middleware.RequireRole(identitydomain.RoleAdmin))

	r.Post("/", h.CreateGalleryPost)
	r.Delete("/{id}", h.DeleteGalleryPost)

	return r
}
