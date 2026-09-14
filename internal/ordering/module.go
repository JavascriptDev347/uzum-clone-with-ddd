package ordering

import (
	cartapp "github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/cart/application"
	catalog "github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/catalog/domain"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/identity/infrastructure/security"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/ordering/application"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/ordering/infrastructure/postgres"
	orderinghttp "github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/ordering/interfaces/http"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

type Module struct {
	// CheckoutRouter va OrdersRouter alohida mount qilinadi (main.go'da) - Catalog'ning
	// router'i allaqachon bare "/api/v1" prefiksini egallagan, chi esa bir xil prefiksni
	// ikki marta mount qilishga ruxsat bermaydi. AdminOrdersRouter mutlaqo yangi prefiksda
	// ("/api/v1/admin/orders") bo'lgani uchun bu cheklovga uchramaydi.
	CheckoutRouter    chi.Router
	OrdersRouter      chi.Router
	AdminOrdersRouter chi.Router

	CheckoutUseCase             *application.CheckoutUseCase
	GetUserOrdersUseCase        *application.GetUserOrdersUseCase
	GetOrderByIDUseCase         *application.GetOrderByIDUseCase
	UpdatePaymentStatusUseCase  *application.UpdatePaymentStatusUseCase
	UpdateDeliveryStatusUseCase *application.UpdateDeliveryStatusUseCase
	GetAllOrdersUseCase         *application.GetAllOrdersUseCase
	CreateManualOrderUseCase    *application.CreateManualOrderUseCase
	CancelOrderUseCase          *application.CancelOrderUseCase

	// HasDeliveredProductUseCase - boshqa context'lar (masalan review) uchun ochilgan,
	// fokuslangan faqat-o'qish so'rovi. Review shu orqali foydalanuvchining mahsulotni
	// sotib olib, yetkazib berilganini tekshiradi - ordering'ning domen yoki repository
	// qatlamiga to'g'ridan-to'g'ri kirmasdan.
	HasDeliveredProductUseCase *application.HasDeliveredProductUseCase
}

type Config struct {
	DB           *sqlx.DB
	TokenService *security.JWTTokenService

	// ProductRepo - Catalog'ning read/write interfeysi: checkout paytida mahsulot nomi va
	// narxini "surat" (snapshot) qilish uchun ishlatiladi.
	ProductRepo catalog.ProductRepository

	// GetCartUseCase va ClearCartUseCase - cart context'ining application-layer chegarasi
	// (ACL): checkout shular orqali cart'ni o'qiydi va bo'shatadi, cart'ning domain yoki
	// repository qatlamiga to'g'ridan-to'g'ri kirmasdan.
	GetCartUseCase   *cartapp.GetCartUseCase
	ClearCartUseCase *cartapp.ClearCartUseCase
}

func NewModule(cfg Config) *Module {
	orderRepo := postgres.NewPostgresOrderRepository(cfg.DB)
	stockReserver := postgres.NewCatalogStockReserver(cfg.ProductRepo)

	checkoutUC := application.NewCheckoutUseCase(orderRepo, cfg.GetCartUseCase, cfg.ClearCartUseCase, cfg.ProductRepo, stockReserver)
	getUserOrdersUC := application.NewGetUserOrdersUseCase(orderRepo)
	getOrderByIDUC := application.NewGetOrderByIDUseCase(orderRepo)
	updatePaymentStatusUC := application.NewUpdatePaymentStatusUseCase(orderRepo)
	updateDeliveryStatusUC := application.NewUpdateDeliveryStatusUseCase(orderRepo)
	getAllOrdersUC := application.NewGetAllOrdersUseCase(orderRepo)
	createManualOrderUC := application.NewCreateManualOrderUseCase(orderRepo, cfg.ProductRepo, stockReserver)
	cancelOrderUC := application.NewCancelOrderUseCase(orderRepo, stockReserver)
	hasDeliveredProductUC := application.NewHasDeliveredProductUseCase(orderRepo)

	orderHandler := orderinghttp.NewOrderHandler(
		checkoutUC,
		getUserOrdersUC,
		getOrderByIDUC,
		updatePaymentStatusUC,
		updateDeliveryStatusUC,
		getAllOrdersUC,
		createManualOrderUC,
		cancelOrderUC,
	)

	return &Module{
		CheckoutRouter:              orderinghttp.NewCheckoutRouter(orderHandler, cfg.TokenService),
		OrdersRouter:                orderinghttp.NewOrdersRouter(orderHandler, cfg.TokenService),
		AdminOrdersRouter:           orderinghttp.NewAdminOrdersRouter(orderHandler, cfg.TokenService),
		CheckoutUseCase:             checkoutUC,
		GetUserOrdersUseCase:        getUserOrdersUC,
		GetOrderByIDUseCase:         getOrderByIDUC,
		UpdatePaymentStatusUseCase:  updatePaymentStatusUC,
		UpdateDeliveryStatusUseCase: updateDeliveryStatusUC,
		GetAllOrdersUseCase:         getAllOrdersUC,
		CreateManualOrderUseCase:    createManualOrderUC,
		CancelOrderUseCase:          cancelOrderUC,
		HasDeliveredProductUseCase:  hasDeliveredProductUC,
	}
}
