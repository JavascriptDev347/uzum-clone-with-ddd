package identity

import (
	"time"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/identity/application"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/identity/infrastructure/postgres"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/identity/infrastructure/security"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/identity/interfaces/http"
	identityhttp "github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/identity/interfaces/http"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

type Module struct {
	Router chi.Router
}

type Config struct {
	DB         *sqlx.DB
	JWTSecret  string
	AccessTTL  time.Duration
	RefreshTTL time.Duration
}

func NewModule(cfg Config) *Module {
	// Infrastructure
	userRepo := postgres.NewPostgresUserRepository(cfg.DB)
	hasher := security.NewBcryptHasher()
	tokenService := security.NewJWTTokenService(cfg.JWTSecret, cfg.AccessTTL, cfg.RefreshTTL)

	// Application
	registerUC := application.NewRegisterUserUseCase(userRepo, hasher)
	loginUC := application.NewLoginUserUseCase(userRepo, hasher, tokenService)
	refreshUC := application.NewRefreshTokenUseCase(userRepo, tokenService)

	// Interfaces
	handler := http.NewIdentityHandler(registerUC, loginUC, refreshUC)

	return &Module{
		Router: identityhttp.NewRouter(handler),
	}
}
