package wishlist

import (
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/identity/infrastructure/security"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/wishlist/application"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/wishlist/infrastructure/postgres"
	wishlisthttp "github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/wishlist/interfaces/http"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

type Module struct {
	Router chi.Router
}

type Config struct {
	DB           *sqlx.DB
	TokenService *security.JWTTokenService
}

func NewModule(cfg Config) *Module {
	wishlistRepo := postgres.NewPostgresWishlistRepository(cfg.DB.DB)

	addUC := application.NewAddToWishlistUseCase(wishlistRepo)
	removeUC := application.NewRemoveFromWishlistUseCase(wishlistRepo)
	getUC := application.NewGetWishlistUseCase(wishlistRepo)

	wishlistHandler := wishlisthttp.NewWishlistHandler(addUC, removeUC, getUC)

	return &Module{
		Router: wishlisthttp.NewRouter(wishlistHandler, cfg.TokenService),
	}
}
