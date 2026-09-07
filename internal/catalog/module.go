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
	createProductUC := application.NewCreateProductUseCase(productRepo, categoryRepo, cfg.MediaUploader)
	getProductsUC := application.NewGetProductsUseCase(productRepo)
	getProductByIDUC := application.NewGetProductUseCase(productRepo)
	getProductBySlugUC := application.NewGetProductBySlugUseCase(productRepo)
	updateProductUC := application.NewUpdateProductUseCase(productRepo, categoryRepo)
	addProductImagesUC := application.NewAddProductImagesUseCase(productRepo, cfg.MediaUploader)
	replaceProductImageUC := application.NewReplaceProductImageUseCase(productRepo, cfg.MediaUploader)
	deleteProductUC := application.NewDeleteProductUseCase(productRepo)
	getAllProductsUC := application.NewGetAllProductsIncludingDeletedUseCase(productRepo)

	createCategoryUC := application.NewCreateCategoryUseCase(categoryRepo, cfg.MediaUploader)
	getCategoriesUC := application.NewGetCategoriesUseCase(categoryRepo)
	getCategoryByIDUC := application.NewGetCategoryUseCase(categoryRepo)
	updateCategoryUC := application.NewUpdateCategoryUseCase(categoryRepo)
	updateCategoryImageUC := application.NewUpdateCategoryImageUseCase(categoryRepo, cfg.MediaUploader)
	deleteCategoryUC := application.NewDeleteCategoryUseCase(categoryRepo)
	getAllCategoriesUC := application.NewGetAllCategoriesIncludingDeletedUseCase(categoryRepo)

	createEventUC := application.NewCreateEventUseCase(eventRepo, categoryRepo, cfg.MediaUploader)
	getEventsUC := application.NewGetEventsUseCase(eventRepo)
	getEventByIDUC := application.NewGetEventUseCase(eventRepo)
	updateEventUC := application.NewUpdateEventUseCase(eventRepo, categoryRepo)
	updateEventImageUC := application.NewUpdateEventImageUseCase(eventRepo, cfg.MediaUploader)
	deleteEventUC := application.NewDeleteEventUseCase(eventRepo)
	getAllEventsUC := application.NewGetAllEventsIncludingDeletedUseCase(eventRepo)

	// interfaces
	productHandler := http.NewProductHandler(createProductUC, getProductsUC, getProductByIDUC, getProductBySlugUC, updateProductUC, addProductImagesUC, replaceProductImageUC, deleteProductUC, getAllProductsUC)
	categoryHandler := http.NewCategoryHandler(createCategoryUC, getCategoriesUC, getCategoryByIDUC, updateCategoryUC, updateCategoryImageUC, deleteCategoryUC, getAllCategoriesUC)
	eventHandler := http.NewEventHandler(createEventUC, getEventsUC, getEventByIDUC, updateEventUC, updateEventImageUC, deleteEventUC, getAllEventsUC)

	return &Module{
		Router: producthttp.NewRouter(productHandler, categoryHandler, eventHandler, cfg.TokenService),
	}
}
