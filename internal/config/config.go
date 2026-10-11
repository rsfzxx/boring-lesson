package config

import (
	"log"
	"os"
	"time"
)

type Config struct {
	Addr         string
	MongoURI     string
	MongoDB      string
	RedisAddr    string
	AMQPURL      string
	JWTSecret    string
	JWTExpiresIn time.Duration
}

func Load() Config {
	cfg := Config{
		Addr:         env("ADDR", ":8080"),
		MongoURI:     mustEnv("MONGO_URI"),
		MongoDB:      mustEnv("MONGO_DB"),
		RedisAddr:    mustEnv("REDIS_ADDR"),
		AMQPURL:      mustEnv("AMQP_URL"),
		JWTSecret:    mustEnv("JWT_SECRET"),
		JWTExpiresIn: duration("JWT_EXPIRES_IN", 24*time.Hour),
	}
	if len(cfg.JWTSecret) < 32 {
		log.Fatal("JWT_SECRET minimal 32 karakter")
	}
	return cfg
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("environment variable %s wajib diisi (lihat .env.example)", key)
	}
	return v
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func duration(key string, fallback time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	d, err := time.ParseDuration(v)
	if err != nil || d <= 0 {
		log.Fatalf("environment variable %s tidak valid, contoh: 24h atau 30m", key)
	}
	return d
}