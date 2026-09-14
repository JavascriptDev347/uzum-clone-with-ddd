package http

import (
	identitydomain "github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/identity/domain"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/identity/infrastructure/security"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/identity/interfaces/http/middleware"
	"github.com/go-chi/chi/v5"
)

// NewReviewsRouter - POST /api/v1/reviews uchun router (autentifikatsiya talab qilinadi).
func NewReviewsRouter(h *ReviewHandler, tokenService *security.JWTTokenService) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.Authenticate(tokenService))

	r.Post("/", h.SubmitReview)

	return r
}

// NewProductReviewsRouter - GET /api/v1/products/{id}/reviews uchun router. Ochiq - har kim
// mahsulot sharhlarini ko'rishi mumkin, autentifikatsiya shart emas. Catalog'ning router'i
// allaqachon bare "/api/v1" prefiksini egallagani uchun (ordering'da ham xuddi shunday
// muammo), bu alohida, aniq prefiks ("/api/v1/products/{id}/reviews") bilan mount qilinadi.
func NewProductReviewsRouter(h *ReviewHandler) chi.Router {
	r := chi.NewRouter()

	r.Get("/", h.GetProductReviews)

	return r
}

// NewAdminReviewsRouter - DELETE /api/v1/admin/reviews/{id} uchun router. Faqat admin -
// istalgan foydalanuvchi yozgan sharhni moderatsiya maqsadida butunlay o'chira oladi.
// "/api/v1/admin/reviews" bare "/api/v1"dan farqli, aniq prefiks bo'lgani uchun
// (gallery/ordering'da bo'lgani kabi) alohida mount qilinadi.
func NewAdminReviewsRouter(h *ReviewHandler, tokenService *security.JWTTokenService) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.Authenticate(tokenService))
	r.Use(middleware.RequireRole(identitydomain.RoleAdmin))

	r.Delete("/{id}", h.DeleteReview)

	return r
}
