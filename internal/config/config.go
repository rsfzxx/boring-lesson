package config

import (
	"log"
	"os"
)

type Config struct {
	Addr      string
	MongoURI  string
	MongoDB   string
	RedisAddr string
	AMQPURL   string
}

func Load() Config {
	return Config{
		Addr:      env("ADDR", ":8080"),
		MongoURI:  mustEnv("MONGO_URI"),
		MongoDB:   mustEnv("MONGO_DB"),
		RedisAddr: mustEnv("REDIS_ADDR"),
		AMQPURL:   mustEnv("AMQP_URL"),
	}
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
