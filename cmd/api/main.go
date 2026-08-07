package main

import (
	"fmt"
	"log"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/pkg/config"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/pkg/database"
)

func main() {
	// config env variables
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	// load db
	db, err := database.NewPostgresDB(cfg.DSN())
	if err != nil {
		log.Fatal(err)
	}
	// load redis
	rdb, err := database.NewRedisClient(cfg.Redis.Host, cfg.Redis.Port)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("db: %v, rdb: %v\n", db, rdb)
	fmt.Println("BRR. Project started !!!")
}
