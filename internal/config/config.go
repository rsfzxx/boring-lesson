// Package config membaca konfigurasi aplikasi dari environment variable.
package config

import "os"

type Config struct {
	Addr      string
	MongoURI  string
	RedisAddr string
	AMQPURL   string
}

func Load() Config {
	return Config{
		Addr:      env("ADDR", ":8080"),
		MongoURI:  env("MONGO_URI", "mongodb://localhost:27017"),
		RedisAddr: env("REDIS_ADDR", "localhost:6379"),
		AMQPURL:   env("AMQP_URL", "amqp://guest:guest@localhost:5672/"),
	}
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
