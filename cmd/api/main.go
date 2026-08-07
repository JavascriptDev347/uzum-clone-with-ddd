package main

import (
	"fmt"
	"log"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/pkg/config"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/pkg/database"
	"github.com/bakhod1r/oneenv"
)

func main() {
	// config env variables
	cfg, err := oneenv.Parse[config.Config]()
	if err != nil {
		log.Fatal(err)
	}

	// load db
	db, err := database.NewPostgresDB(cfg.DSN())
	if err != nil {
		log.Fatal(err)
	}
	// load redis
	rdb, err := database.NewRedisClient(cfg.RedisHost, cfg.RedisPort)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("db: %v, rdb: %v\n", db, rdb)
	fmt.Println("BRR. Project started !!!")
}
