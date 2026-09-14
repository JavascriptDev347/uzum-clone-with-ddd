package review

import (
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/identity/infrastructure/security"
	orderingapp "github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/ordering/application"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/review/application"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/review/infrastructure/postgres"
	reviewhttp "github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/review/interfaces/http"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

type Module struct {
	// ReviewsRouter va ProductReviewsRouter alohida mount qilinadi (main.go'da) - Catalog'ning
	// router'i allaqachon bare "/api/v1" prefiksini egallagan (ordering'da ham xuddi shunday
	// muammo hal qilingan edi).
	ReviewsRouter        chi.Router
	ProductReviewsRouter chi.Router

	SubmitReviewUseCase      *application.SubmitReviewUseCase
	GetProductReviewsUseCase *application.GetProductReviewsUseCase
	GetUserReviewsUseCase    *application.GetUserReviewsUseCase
}

type Config struct {
	DB           *sqlx.DB
	TokenService *security.JWTTokenService

	// HasDeliveredProductUseCase - ordering context'ining application-layer chegarasi (ACL):
	// review shu orqali foydalanuvchining mahsulotni sotib olib, yetkazib berilganini
	// tekshiradi, ordering'ning domain yoki repository qatlamiga to'g'ridan-to'g'ri
	// kirmasdan (cart va catalog o'rtasidagi hozirgi to'g'ridan-to'g'ri bog'lanishdan farqli
	// o'laroq, bu context uchun ACL chegarasi qat'iy saqlanadi).
	HasDeliveredProductUseCase *orderingapp.HasDeliveredProductUseCase
}

func NewModule(cfg Config) *Module {
	reviewRepo := postgres.NewPostgresReviewRepository(cfg.DB)
	purchaseChecker := postgres.NewOrderingPurchaseChecker(cfg.HasDeliveredProductUseCase)

	submitReviewUC := application.NewSubmitReviewUseCase(reviewRepo, purchaseChecker)
	getProductReviewsUC := application.NewGetProductReviewsUseCase(reviewRepo)
	getUserReviewsUC := application.NewGetUserReviewsUseCase(reviewRepo)

	reviewHandler := reviewhttp.NewReviewHandler(submitReviewUC, getProductReviewsUC, getUserReviewsUC)

	return &Module{
		ReviewsRouter:            reviewhttp.NewReviewsRouter(reviewHandler, cfg.TokenService),
		ProductReviewsRouter:     reviewhttp.NewProductReviewsRouter(reviewHandler),
		SubmitReviewUseCase:      submitReviewUC,
		GetProductReviewsUseCase: getProductReviewsUC,
		GetUserReviewsUseCase:    getUserReviewsUC,
	}
}
