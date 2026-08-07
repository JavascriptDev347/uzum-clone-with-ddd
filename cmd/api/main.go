package main

import (
	"log"
	"net/http"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/identity"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/pkg/config"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/pkg/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

	// standart middlewares
	r := chi.NewRouter()
	// //X-Request-Id ga uuiddan biriktirib qaytaradi
	r.Use(middleware.RequestID)
	//haqiqiy IP-manzilini aniqlaydi.
	r.Use(middleware.RealIP)
	// log qilib yozadi hammasini
	r.Use(middleware.Logger)
	//Kodda kutilmagan xatolik tufayli panic() yuz berganda serverning to'xtab (crash bo'lib) qolishining oldini oladi.
	r.Use(middleware.Recoverer)

	r.Mount("/api/v1/auth", identityModule.Router)
	// ── Server ────────────────────────────────────────────────
	log.Printf("server started on :%s", cfg.AppPort)
	if err := http.ListenAndServe(":8080", identityModule.Router); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
