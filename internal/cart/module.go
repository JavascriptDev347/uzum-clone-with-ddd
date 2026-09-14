package cart

import (
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/cart/application"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/cart/infrastructure/postgres"
	carthttp "github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/cart/interfaces/http"
	catalogpostgres "github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/catalog/infrastructure/postgres"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/identity/infrastructure/security"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

type Module struct {
	Router chi.Router

	// GetCartUseCase va ClearCartUseCase boshqa context'lar (masalan ordering) uchun
	// ochiq qilingan - bu cart'ning application-layer chegarasi (ACL): tashqi context'lar
	// cart'ning domain yoki repository qatlamiga to'g'ridan-to'g'ri kirmasdan shu orqali ishlaydi.
	GetCartUseCase   *application.GetCartUseCase
	ClearCartUseCase *application.ClearCartUseCase
}

type Config struct {
	DB           *sqlx.DB
	TokenService *security.JWTTokenService
}

func NewModule(cfg Config) *Module {
	cartRepo := postgres.NewPostgresCartRepository(cfg.DB.DB)
	productRepo := catalogpostgres.NewPostgresProductRepository(cfg.DB)

	addItemUC := application.NewAddItemUseCase(cartRepo, productRepo)
	removeItemUC := application.NewRemoveItemUseCase(cartRepo)
	updateItemQuantityUC := application.NewUpdateItemQuantityUseCase(cartRepo, productRepo)
	getCartUC := application.NewGetCartUseCase(cartRepo, productRepo)
	clearCartUC := application.NewClearCartUseCase(cartRepo)

	cartHandler := carthttp.NewCartHandler(addItemUC, removeItemUC, updateItemQuantityUC, getCartUC)

	return &Module{
		Router:           carthttp.NewRouter(cartHandler, cfg.TokenService),
		GetCartUseCase:   getCartUC,
		ClearCartUseCase: clearCartUC,
	}
}
