package main

import (
	"context"
	"log"
	"os"
	"time"

	"go-starter/database/migrations"
	"go-starter/internal/config"
	"go-starter/internal/infra"
)

func main() {
	if len(os.Args) != 2 || (os.Args[1] != "up" && os.Args[1] != "down") {
		log.Fatal("pemakaian: migrate up | migrate down")
	}

	cfg := config.Load()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := infra.NewMongo(ctx, cfg.MongoURI)
	if err != nil {
		log.Fatal("mongodb: ", err)
	}
	defer client.Disconnect(context.Background())

	db := client.Database(cfg.MongoDB)
	if os.Args[1] == "up" {
		err = migrations.Up(ctx, db)
	} else {
		err = migrations.Down(ctx, db)
	}
	if err != nil {
		log.Fatal(err)
	}
	log.Println("selesai")
}