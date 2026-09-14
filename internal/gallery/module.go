package gallery

import (
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/gallery/application"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/gallery/infrastructure/postgres"
	galleryhttp "github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/gallery/interfaces/http"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/identity/infrastructure/security"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/shared/media"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

type Module struct {
	// GalleryRouter va AdminGalleryRouter alohida mount qilinadi (main.go'da) - Catalog'ning
	// router'i allaqachon bare "/api/v1" prefiksini egallagan (ordering/review'da ham xuddi
	// shunday muammo hal qilingan edi).
	GalleryRouter      chi.Router
	AdminGalleryRouter chi.Router

	CreateGalleryPostUseCase *application.CreateGalleryPostUseCase
	ListGalleryPostsUseCase  *application.ListGalleryPostsUseCase
	DeleteGalleryPostUseCase *application.DeleteGalleryPostUseCase
}

type Config struct {
	DB            *sqlx.DB
	TokenService  *security.JWTTokenService
	MediaUploader media.Uploader
}

func NewModule(cfg Config) *Module {
	repo := postgres.NewPostgresGalleryPostRepository(cfg.DB)

	createUC := application.NewCreateGalleryPostUseCase(repo, cfg.MediaUploader)
	listUC := application.NewListGalleryPostsUseCase(repo)
	deleteUC := application.NewDeleteGalleryPostUseCase(repo, cfg.MediaUploader)

	handler := galleryhttp.NewGalleryHandler(createUC, listUC, deleteUC)

	return &Module{
		GalleryRouter:            galleryhttp.NewGalleryRouter(handler),
		AdminGalleryRouter:       galleryhttp.NewAdminGalleryRouter(handler, cfg.TokenService),
		CreateGalleryPostUseCase: createUC,
		ListGalleryPostsUseCase:  listUC,
		DeleteGalleryPostUseCase: deleteUC,
	}
}
