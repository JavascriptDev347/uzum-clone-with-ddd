package http

import (
	identitydomain "github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/identity/domain"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/identity/infrastructure/security"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/identity/interfaces/http/middleware"
	"github.com/go-chi/chi/v5"
)

// NewCheckoutRouter - POST /api/v1/checkout uchun router. Catalog'ning router'i allaqachon
// bare "/api/v1" prefiksini egallagani uchun (chi bir xil prefiksni ikki marta mount qilishga
// ruxsat bermaydi), checkout va orders/{id} alohida-alohida, aniq sub-prefikslar bilan
// (main.go'da) mount qilinadi.
func NewCheckoutRouter(h *OrderHandler, tokenService *security.JWTTokenService) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.Authenticate(tokenService))

	r.Post("/", h.Checkout)

	return r
}

// NewOrdersRouter - GET /api/v1/orders va GET /api/v1/orders/{id} uchun router. Admin-only
// yo'llar (payment-status, delivery-status, admin navbati) ham shu prefiks ostida, chunki
// bare "/api/v1" allaqachon Catalog'ga tegishli va uchinchi router yaratish shart emas.
// "/admin" statik yo'l bo'lgani uchun "/{id}" wildcard'dan oldin ustuvor bo'ladi (chi'ning
// odatiy xulqi - Catalog'ning /products/admin vs /products/{id}'da ham xuddi shunday).
func NewOrdersRouter(h *OrderHandler, tokenService *security.JWTTokenService) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.Authenticate(tokenService))

	r.Get("/", h.GetUserOrders)
	r.Get("/{id}", h.GetOrderByID)

	r.With(middleware.RequireRole(identitydomain.RoleAdmin)).Get("/admin", h.GetAllOrders)
	r.With(middleware.RequireRole(identitydomain.RoleAdmin)).Patch("/{id}/payment-status", h.UpdatePaymentStatus)
	r.With(middleware.RequireRole(identitydomain.RoleAdmin)).Patch("/{id}/delivery-status", h.UpdateDeliveryStatus)

	return r
}

// NewAdminOrdersRouter - POST /api/v1/admin/orders uchun router: admin tomonidan cart'siz,
// to'g'ridan-to'g'ri buyurtma yaratish (offline savdo). "/api/v1/admin/orders" mutlaqo yangi,
// hech kim egallamagan prefiks bo'lgani uchun (Checkout/Orders uchun ishlatilgan aylanma
// yechim shart emas) alohida, oddiy router sifatida qo'shildi.
func NewAdminOrdersRouter(h *OrderHandler, tokenService *security.JWTTokenService) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.Authenticate(tokenService))
	r.Use(middleware.RequireRole(identitydomain.RoleAdmin))

	r.Post("/", h.CreateManualOrder)
	r.Patch("/{id}/cancel", h.CancelOrder)

	return r
}
