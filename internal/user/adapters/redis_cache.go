package adapters

import (
	"context"
	"encoding/json"
	"time"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"

	"scaffold/internal/common/utils"
	"scaffold/internal/user/domain"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(client *redis.Client) domain.CacheRepository {
	return &RedisCache{client: client}
}

const (
	keyRefreshTokenMapDuration = 30 * 24 * time.Hour
	keyRefreshTokenMap         = "refresh_token_map"
)

func (ch *RedisCache) GenRefreshToken(payload *domain.JwtPayload) (string, error) {
	refreshToken, err := utils.GenRandomHexToken()
	if err != nil {
		return "", errors.WithStack(err)
	}

	key := utils.GetRedisKey(keyRefreshTokenMap)
	payloadByte, err := json.Marshal(payload)
	if err != nil {
		return "", errors.WithStack(err)
	}

	payloadStr := string(payloadByte)

	pipe := ch.client.Pipeline()

	if err := pipe.HSet(context.Background(), key, payloadStr, refreshToken).Err(); err != nil {
		return "", errors.WithStack(err)
	}

	// 使用计算出的过期时间
	pipe.HExpire(context.Background(), key, keyRefreshTokenMapDuration, payloadStr)

	// 执行Pipeline命令
	_, err = pipe.Exec(context.Background())
	if err != nil {
		return "", errors.WithStack(err)
	}

	return refreshToken, nil
}

func (ch *RedisCache) ValidateRefreshToken(payload *domain.JwtPayload, refreshToken string) error {
	payloadByte, err := json.Marshal(payload)
	if err != nil {
		return errors.WithStack(err)
	}

	key := utils.GetRedisKey(keyRefreshTokenMap)
	payloadStr := string(payloadByte)

	result, err := ch.client.HGet(context.Background(), key, payloadStr).Result()
	if err != nil {
		if err == redis.Nil {
			return errors.New("refresh token not found or expired")
		}
		return errors.WithStack(err)
	}

	if refreshToken != result {
		return errors.New("refresh token invalid")
	}

	return nil
}

func (ch *RedisCache) ResetRefreshTokenExpiry(payload *domain.JwtPayload) error {
	payloadByte, err := json.Marshal(payload)
	if err != nil {
		return errors.WithStack(err)
	}

	key := utils.GetRedisKey(keyRefreshTokenMap)
	payloadStr := string(payloadByte)

	if err := ch.client.HExpire(context.Background(), key, keyRefreshTokenMapDuration, payloadStr).Err(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
