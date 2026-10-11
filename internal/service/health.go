// Package service berisi logika bisnis, tanpa tahu soal HTTP.
package service

import (
	"context"
	"errors"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"go-starter/internal/dto"
)

type HealthService struct {
	mongo  *mongo.Client
	redis  *redis.Client
	rabbit *amqp.Connection
}

func NewHealthService(m *mongo.Client, r *redis.Client, q *amqp.Connection) *HealthService {
	return &HealthService{mongo: m, redis: r, rabbit: q}
}

// Check memeriksa setiap dependency dan merangkum hasilnya.
func (s *HealthService) Check(ctx context.Context) dto.HealthResponse {
	checks := map[string]error{
		"mongodb":  s.mongo.Ping(ctx, nil),
		"redis":    s.redis.Ping(ctx).Err(),
		"rabbitmq": s.rabbitStatus(),
	}

	res := dto.HealthResponse{Status: "ok", Services: make(map[string]string, len(checks))}
	for name, err := range checks {
		if err != nil {
			res.Status = "degraded"
			res.Services[name] = "down: " + err.Error()
			continue
		}
		res.Services[name] = "up"
	}
	return res
}

func (s *HealthService) rabbitStatus() error {
	if s.rabbit.IsClosed() {
		return errors.New("connection closed")
	}
	return nil
}
