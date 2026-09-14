// @title Uzum Clone API
// @version 1.0
// @description DDD asosida qurilgan e-commerce marketplace backend
// @host
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

package main

import (
	"log"
	"net/http"

	_ "github.com/JavascriptDev347/uzum-clone-with-ddd.git/docs"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/admin/dashboard"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/cart"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/catalog"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/gallery"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/identity"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/ordering"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/review"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/shared/media"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/wishlist"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/pkg/config"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/pkg/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	// config env variables
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	// ── Infrastructure (umumiy, barcha context'lar ishlatadi) ──
	db, err := database.NewPostgresDB(cfg.DSN())
	if err != nil {
		log.Fatalf("postgres connection error: %v", err)
	}
	defer db.Close()

	rdb, err := database.NewRedisClient(cfg.Redis.Host, cfg.Redis.Port)
	if err != nil {
		log.Fatalf("redis connection error: %v", err)
	}
	defer rdb.Close()

	// ── Bounded context modullari ────────────────────────────
	identityModule := identity.NewModule(identity.Config{
		DB:         db,
		JWTSecret:  cfg.JWT.Secret,
		AccessTTL:  cfg.JWT.AccessTTL,
		RefreshTTL: cfg.JWT.RefreshTTL,
	})

	s3Uploader, err := media.NewS3Uploader(media.S3Config{
		AccessKeyID:     cfg.AWS.AccessKeyID,
		SecretAccessKey: cfg.AWS.SecretAccessKey,
		Region:          cfg.AWS.Region,
		Bucket:          cfg.AWS.S3BucketName,
		Folder:          cfg.AWS.S3Folder,
	})
	if err != nil {
		log.Fatalf("s3 init error: %v", err)
	}

	catalogModule := catalog.NewModule(catalog.Config{
		DB:            db,
		TokenService:  identityModule.TokenService,
		MediaUploader: s3Uploader,
	})

	wishlistModule := wishlist.NewModule(wishlist.Config{
		DB:           db,
		TokenService: identityModule.TokenService,
	})

	cartModule := cart.NewModule(cart.Config{
		DB:           db,
		TokenService: identityModule.TokenService,
	})

	orderingModule := ordering.NewModule(ordering.Config{
		DB:               db,
		TokenService:     identityModule.TokenService,
		ProductRepo:      catalogModule.ProductRepo,
		GetCartUseCase:   cartModule.GetCartUseCase,
		ClearCartUseCase: cartModule.ClearCartUseCase,
	})

	reviewModule := review.NewModule(review.Config{
		DB:                         db,
		TokenService:               identityModule.TokenService,
		HasDeliveredProductUseCase: orderingModule.HasDeliveredProductUseCase,
	})

	galleryModule := gallery.NewModule(gallery.Config{
		DB:            db,
		TokenService:  identityModule.TokenService,
		MediaUploader: s3Uploader,
	})

	// standart middlewares
	r := chi.NewRouter()
	// cors
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"}, // dev uchun barchasi, production'da aniq domainlar yozish kerak
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Requested-With"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false, // agar cookie/credentials ishlatilsa true qil, lekin unda AllowedOrigins "*" bo'la olmaydi
		MaxAge:           300,
	}))
	// //X-Request-Id ga uuiddan biriktirib qaytaradi
	r.Use(middleware.RequestID)
	//haqiqiy IP-manzilini aniqlaydi.
	r.Use(middleware.RealIP)
	// log qilib yozadi hammasini
	r.Use(middleware.Logger)
	//Kodda kutilmagan xatolik tufayli panic() yuz berganda serverning to'xtab (crash bo'lib) qolishining oldini oladi.
	r.Use(middleware.Recoverer)

	r.Mount("/api/v1/auth", identityModule.Router)
	r.Mount("/api/v1", catalogModule.Router)
	r.Mount("/api/v1/wishlist", wishlistModule.Router)
	r.Mount("/api/v1/cart", cartModule.Router)
	// Catalog'ning router'i bare "/api/v1"ni egallagani uchun (chi bir xil prefiksni ikki
	// marta mount qilishga ruxsat bermaydi), checkout va orders alohida, aniq prefikslar bilan
	// mount qilinadi.
	r.Mount("/api/v1/checkout", orderingModule.CheckoutRouter)
	r.Mount("/api/v1/orders", orderingModule.OrdersRouter)
	r.Mount("/api/v1/admin/orders", orderingModule.AdminOrdersRouter)
	// Xuddi shu sabab bilan (Catalog bare "/api/v1"ni egallagan) - review'ning ikkala router'i
	// ham aniq, mos ravishda "/api/v1/products"dan ham ustuvor bo'ladigan ("/{id}/reviews"
	// bilan tugaydigan, chi'ning radix daraxti bo'yicha aniqroq mos kelish) prefikslar bilan
	// mount qilinadi.
	r.Mount("/api/v1/reviews", reviewModule.ReviewsRouter)
	r.Mount("/api/v1/products/{id}/reviews", reviewModule.ProductReviewsRouter)
	// "/api/v1/gallery" va "/api/v1/admin/gallery" ham bare "/api/v1"dan farqli, aniq
	// prefikslar - Catalog bilan to'qnashmaydi.
	r.Mount("/api/v1/gallery", galleryModule.GalleryRouter)
	r.Mount("/api/v1/admin/gallery", galleryModule.AdminGalleryRouter)
	// admin/dashboard - bounded context emas, module.go yo'q, to'g'ridan-to'g'ri *sqlx.DB'ga
	// tayanadi (internal/admin/dashboard/router.go'ga qarang).
	r.Mount("/api/v1/admin/dashboard", dashboard.NewRouter(db, identityModule.TokenService))

	// ── Swagger UI: http://localhost:8080/swagger/index.html ──
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	// ── Server ────────────────────────────────────────────────
	log.Printf("server started on :%d", cfg.AppPort)
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
