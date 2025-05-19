package infrastructure

import (
	"context"
	"encoding/json"
	"scaffold/internal/domain/user/model"
	"scaffold/utils"
	"time"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

type ICache interface {
	GenRefreshToken(payload *model.JwtPayload) (string, error)
	ValidateRefreshToken(payload *model.JwtPayload, refreshToken string) error
	ResetRefreshTokenExpiry(payload *model.JwtPayload) error
}

type cache struct {
	client *redis.Client
}

func NewCache(client *redis.Client) ICache {
	return &cache{client: client}
}

const (
	keyRefreshTokenMapDuration = 30 * 24 * time.Hour
	keyRefreshTokenMap         = "refresh_token_map" //这里使用map只是为了以后方便复用 该项目实际只需要一个string键就可以了
)

func (ch *cache) GenRefreshToken(payload *model.JwtPayload) (string, error) {
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

func (ch *cache) ValidateRefreshToken(payload *model.JwtPayload, refreshToken string) error {
	payloadByte, err := json.Marshal(payload)
	if err != nil {
		return errors.WithStack(err)
	}
	key := utils.GetRedisKey(keyRefreshTokenMap)

	payloadStr := string(payloadByte)

	result, err := ch.client.HGet(context.Background(), key, payloadStr).Result()
	if err != nil {
		return errors.WithStack(err)
	}

	if refreshToken == result {
		return nil
	}
	return errors.New("refreshToken 无效")
}

func (ch *cache) ResetRefreshTokenExpiry(payload *model.JwtPayload) error {
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
