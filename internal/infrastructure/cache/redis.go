package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/smokersaim/droqsic/cmd/config"
	"go.uber.org/zap"
)

type RedisCache struct {
	Client *redis.Client
}

func ConnectRedisCache(cfg *config.Config, log *zap.Logger) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Cache.URI,
		DB:           cfg.Cache.Database,
		Username:     cfg.Cache.Username,
		Password:     cfg.Cache.Password,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     50,
		MinIdleConns: 10,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		log.Error("Redis connection failed", zap.Error(err))
		return nil, err
	}

	log.Info("Redis cache initialized successfully")
	return &RedisCache{Client: client}, nil
}

func (r *RedisCache) Disconnect(log *zap.Logger) {
	if err := r.Client.Close(); err != nil {
		log.Error("Error while disconnecting from Redis", zap.Error(err))
	} else {
		log.Info("Redis cache disconnected successfully")
	}
}
