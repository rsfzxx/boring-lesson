package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-starter/internal/config"
	"go-starter/internal/handler"
	"go-starter/internal/infra"
	"go-starter/internal/router"
	"go-starter/internal/service"
)

func main() {
	cfg := config.Load()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	startCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	mongoClient, err := infra.NewMongo(startCtx, cfg.MongoURI)
	if err != nil {
		log.Fatal("mongodb: ", err)
	}
	redisClient, err := infra.NewRedis(startCtx, cfg.RedisAddr)
	if err != nil {
		log.Fatal("redis: ", err)
	}
	rabbitConn, err := infra.NewRabbitMQ(cfg.AMQPURL)
	if err != nil {
		log.Fatal("rabbitmq: ", err)
	}

	healthSvc := service.NewHealthService(mongoClient, redisClient, rabbitConn)
	healthHandler := handler.NewHealthHandler(healthSvc)

	srv := &http.Server{Addr: cfg.Addr, Handler: router.New(healthHandler)}

	go func() {
		log.Println("API jalan di", cfg.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("server: ", err)
		}
	}()

	<-ctx.Done()
	log.Println("mematikan server...")

	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelShutdown()
	srv.Shutdown(shutdownCtx)
	rabbitConn.Close()
	redisClient.Close()
	mongoClient.Disconnect(shutdownCtx)
}
