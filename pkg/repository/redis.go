package repository

import (
	"context"
	"fmt"
	"scaffold/pkg/config"
	"time"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var client *redis.Client

func RedisInit(cfg *config.RedisConfig) error {
	client = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s",
			cfg.Host,
			cfg.Port,
		),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	if _, err := client.Ping(ctx).Result(); err != nil {
		return errors.Wrap(err, "Redis连接失败")
	}

	zap.L().Info("Redis连接成功",
		zap.String("host", cfg.Host))

	return nil
}

func RedisClose() {
	if err := client.Close(); err != nil {
		zap.L().Error("redis 关闭连接失败", zap.Error(err))
	}
}

func RedisClient() *redis.Client {
	return client
}
