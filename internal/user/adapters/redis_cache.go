package adapters

import (
	"context"
	"encoding/json"
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"

	"scaffold/internal/common/utils"
	"scaffold/internal/user/domain"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisTokenCache() domain.TokenCache {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	password := os.Getenv("REDIS_PASSWORD")
	dbStr := os.Getenv("REDIS_DB")
	poolSizeStr := os.Getenv("REDIS_POOL_SIZE")

	db, _ := strconv.Atoi(dbStr)
	poolSize, _ := strconv.Atoi(poolSizeStr)

	addr := host + ":" + port

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		DB:       db,
		Password: password,
		PoolSize: poolSize,
	})

	// 可选：ping 检查连接
	if err := client.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}

	return &RedisCache{client: client}
}

const (
	keyRefreshTokenMapDuration = 30 * 24 * time.Hour
	keyRefreshTokenMap         = "user_refresh_token_map"
)

func (ch *RedisCache) GenRefreshToken(payload domain.JwtPayload) (string, error) {
	refreshToken, err := utils.GenRandomHexToken()
	if err != nil {
		return "", errors.WithStack(err)
	}

	key := utils.GetRedisKey(keyRefreshTokenMap)
	pipe := ch.client.Pipeline()

	payloadByte, err := json.Marshal(payload)
	if err != nil {
		return "", errors.WithStack(err)
	}
	payloadStr := string(payloadByte)

	if err := pipe.HSet(context.Background(), key, payloadStr, refreshToken).Err(); err != nil {
		return "", errors.WithStack(err)
	}

	pipe.HExpire(context.Background(), key, keyRefreshTokenMapDuration, payloadStr)

	// 执行Pipeline命令
	_, err = pipe.Exec(context.Background())
	if err != nil {
		return "", errors.WithStack(err)
	}

	return refreshToken, nil
}

func (ch *RedisCache) ValidateRefreshToken(payload domain.JwtPayload, refreshToken string) error {
	key := utils.GetRedisKey(keyRefreshTokenMap)

	payloadByte, err := json.Marshal(payload)
	if err != nil {
		return errors.WithStack(err)
	}
	payloadStr := string(payloadByte)

	result, err := ch.client.HGet(context.Background(), key, payloadStr).Result()

	if err != nil {
		if errors.Is(err, redis.Nil) {
			return errors.New("refresh token not found or expired")
		}
		return errors.WithStack(err)
	}

	if refreshToken != result {
		return errors.New("refresh token invalid")
	}

	return nil
}

func (ch *RedisCache) ResetRefreshTokenExpiry(payload domain.JwtPayload) error {
	key := utils.GetRedisKey(keyRefreshTokenMap)

	payloadByte, err := json.Marshal(payload)
	if err != nil {
		return errors.WithStack(err)
	}
	payloadStr := string(payloadByte)

	if err := ch.client.HExpire(context.Background(), key, keyRefreshTokenMapDuration, payloadStr).Err(); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
