package adapters

import (
	"context"
	"encoding/json"
	"scaffold/internal/common/reskit/codes"
	"time"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"

	"scaffold/internal/common/utils"
	"scaffold/internal/user/domain"
)

type TokenRedisCache struct {
	client *redis.Client
}

func NewTokenRedisCache() domain.TokenCache {
	host := utils.GetEnv("REDIS_HOST")
	port := utils.GetEnv("REDIS_PORT")
	password := utils.GetEnv("REDIS_PASSWORD")
	db := utils.GetEnvAsInt("REDIS_DB")
	poolSize := utils.GetEnvAsInt("REDIS_POOL_SIZE")

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

	return &TokenRedisCache{client: client}
}

const (
	keyRefreshTokenMapDuration = 30 * 24 * time.Hour
	keyRefreshTokenMap         = "user_refresh_token_map"
)

func (ch *TokenRedisCache) GenRefreshToken(payload *domain.JwtPayload) (string, error) {
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

	if err := pipe.HSet(context.Background(), key, refreshToken, payloadStr).Err(); err != nil {
		return "", errors.WithStack(err)
	}

	pipe.HExpire(context.Background(), key, keyRefreshTokenMapDuration, refreshToken)

	// 执行Pipeline命令
	_, err = pipe.Exec(context.Background())
	if err != nil {
		return "", errors.WithStack(err)
	}

	return refreshToken, nil
}

func (ch *TokenRedisCache) ValidateRefreshToken(refreshToken string) (*domain.JwtPayload, error) {
	key := utils.GetRedisKey(keyRefreshTokenMap)

	result, err := ch.client.HGet(context.Background(), key, refreshToken).Result()

	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, codes.ErrRefreshTokenNotFound
		}
		return nil, errors.WithStack(err)
	}

	payload := new(domain.JwtPayload)
	if err := json.Unmarshal([]byte(result), payload); err != nil {
		return nil, err
	}

	return payload, nil
}

func (ch *TokenRedisCache) RemoveRefreshToken(refreshToken string) error {
	key := utils.GetRedisKey(keyRefreshTokenMap)

	if err := ch.client.HDel(context.Background(), key, refreshToken).Err(); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
