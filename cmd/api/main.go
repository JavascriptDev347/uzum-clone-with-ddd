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
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/catalog"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/identity"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/shared/media"
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

	cloudinaryUploader, err := media.NewCloudinaryUploader(media.CloudinaryConfig{
		CloudName: cfg.Cloudinary.CloudName,
		APIKey:    cfg.Cloudinary.APIKey,
		APISecret: cfg.Cloudinary.APISecret,
		Folder:    cfg.Cloudinary.Folder,
	})
	if err != nil {
		log.Fatalf("cloudinary init error: %v", err)
	}

	catalogModule := catalog.NewModule(catalog.Config{
		DB:            db,
		TokenService:  identityModule.TokenService,
		MediaUploader: cloudinaryUploader,
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

	// ── Swagger UI: http://localhost:8080/swagger/index.html ──
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	// ── Server ────────────────────────────────────────────────
	log.Printf("server started on :%s", cfg.AppPort)
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
