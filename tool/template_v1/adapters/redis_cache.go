package adapters

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"os"
	"scaffold/internal/common/utils"
	"scaffold/internal/user/model"
	"strconv"
	"time"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache() *RedisCache {
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
	keyRefreshTokenMap         = "refresh_token_map" //这里使用map只是为了以后方便复用 该项目实际只需要一个string键就可以了
)

func (ch *RedisCache) GenRefreshToken(payload *model.JwtPayload) (string, error) {
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

func (ch *RedisCache) ValidateRefreshToken(payload *model.JwtPayload, refreshToken string) error {
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

func (ch *RedisCache) ResetRefreshTokenExpiry(payload *model.JwtPayload) error {
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
