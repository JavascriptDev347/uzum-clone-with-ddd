package dashboard

import (
	identitydomain "github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/identity/domain"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/identity/infrastructure/security"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/identity/interfaces/http/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

// NewRouter - GET /api/v1/admin/dashboard/summary, .../revenue-history va .../low-stock
// uchun router. Bu qatlam bounded context emas (domain/application/infrastructure
// ajratilmagan), shuning uchun alohida module.go yo'q - to'g'ridan-to'g'ri *sqlx.DB'ga
// tayanadi. Uchalasi ham faqat admin uchun.
func NewRouter(db *sqlx.DB, tokenService *security.JWTTokenService) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.Authenticate(tokenService))
	r.Use(middleware.RequireRole(identitydomain.RoleAdmin))

	r.Get("/summary", SummaryHandler(db))
	r.Get("/revenue-history", RevenueHistoryHandler(db))
	r.Get("/low-stock", LowStockHandler(db))

	return r
}
