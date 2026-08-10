package catalog

import (
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/catalog/application"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/catalog/infrastructure/postgres"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/catalog/interfaces/http"
	producthttp "github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/catalog/interfaces/http"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

type Module struct {
	Router chi.Router
}

type Config struct {
	DB *sqlx.DB
}

func NewModule(cfg Config) *Module {
	productRepo := postgres.NewPostgresProductRepository(cfg.DB)

	// application
	createProductUC := application.NewCreateProductUseCase(productRepo)

	// interfaces
	handler := http.NewProductHandler(createProductUC)

	return &Module{
		Router: producthttp.NewRouter(handler),
	}
}
