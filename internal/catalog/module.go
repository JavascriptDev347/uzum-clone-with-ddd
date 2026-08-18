package catalog

import (
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/catalog/application"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/catalog/infrastructure/postgres"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/catalog/interfaces/http"
	producthttp "github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/catalog/interfaces/http"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/identity/infrastructure/security"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/shared/media"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

type Module struct {
	Router chi.Router
}

type Config struct {
	DB            *sqlx.DB
	TokenService  *security.JWTTokenService
	MediaUploader media.Uploader
}

func NewModule(cfg Config) *Module {

	productRepo := postgres.NewPostgresProductRepository(cfg.DB)
	categoryRepo := postgres.NewPostgresCategoryRepository(cfg.DB)
	eventRepo := postgres.NewPostgresEventRepository(cfg.DB)

	// application
	createProductUC := application.NewCreateProductUseCase(productRepo, cfg.MediaUploader)
	createCategoryUC := application.NewCreateCategoryUseCase(categoryRepo, cfg.MediaUploader)
	getCategoriesUC := application.NewGetCategoriesUseCase(categoryRepo)
	getCategoryByIDUC := application.NewGetCategoryUseCase(categoryRepo)
	updateCategoryUC := application.NewUpdateCategoryUseCase(categoryRepo)
	deleteCategoryUC := application.NewDeleteCategoryUseCase(categoryRepo)
	getAllCategoriesUC := application.NewGetAllCategoriesIncludingDeletedUseCase(categoryRepo)

	createEventUC := application.NewCreateEventUseCase(eventRepo, categoryRepo, cfg.MediaUploader)
	getEventsUC := application.NewGetEventsUseCase(eventRepo)
	getEventByIDUC := application.NewGetEventUseCase(eventRepo)
	updateEventUC := application.NewUpdateEventUseCase(eventRepo, categoryRepo)
	deleteEventUC := application.NewDeleteEventUseCase(eventRepo)
	getAllEventsUC := application.NewGetAllEventsIncludingDeletedUseCase(eventRepo)

	// interfaces
	productHandler := http.NewProductHandler(createProductUC)
	categoryHandler := http.NewCategoryHandler(createCategoryUC, getCategoriesUC, getCategoryByIDUC, updateCategoryUC, deleteCategoryUC, getAllCategoriesUC)
	eventHandler := http.NewEventHandler(createEventUC, getEventsUC, getEventByIDUC, updateEventUC, deleteEventUC, getAllEventsUC)

	return &Module{
		Router: producthttp.NewRouter(productHandler, categoryHandler, eventHandler, cfg.TokenService),
	}
}
