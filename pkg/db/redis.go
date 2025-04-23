package db

import (
	"context"
	"fmt"
	"time"

	"bet-settlement-engine/pkg/env"

	"github.com/go-redis/redis/v8"
)

// RedisConfig holds Redis environment variables.
type RedisConfig struct {
	Host     string `env:"REDIS_HOST"`
	Port     string `env:"REDIS_PORT"`
	Password string `env:"REDIS_PASSWORD"`
	DB       int    `env:"REDIS_DB"`
}

// RedisDB wraps the Redis client and its configuration.
type RedisDB struct {
	conf *RedisConfig
	rdb  *redis.Client
}

// Init initializes the Redis client.
func (r *RedisDB) Init() error {
	var cfg RedisConfig
	if err := env.Parse(&cfg); err != nil {
		return fmt.Errorf("failed to parse Redis config: %w", err)
	}
	r.conf = &cfg

	addr := fmt.Sprintf("%s:%s", r.conf.Host, r.conf.Port)
	r.rdb = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: r.conf.Password,
		DB:       r.conf.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := r.rdb.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("unable to connect to Redis at %s: %w", addr, err)
	}

	// Replace with your logger if needed
	fmt.Println("✅ Connected to Redis:", addr)

	return nil
}

// Stop closes the Redis client connection.
func (r *RedisDB) Stop() error {
	if r.rdb != nil {
		return r.rdb.Close()
	}
	return nil
}

// GetClient returns the active Redis client.
func (r *RedisDB) GetClient() *redis.Client {
	return r.rdb
}
