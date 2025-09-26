package adapters

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"os"
	"scaffold/internal/captcha/domain"
	"scaffold/internal/common/reskit/codes"
	"scaffold/internal/common/uid"
	"scaffold/internal/common/utils"
	"strconv"
)

type RedisCaptchaCache struct {
	client *redis.Client
}

func NewRedisCache() domain.CaptchaCache {
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

	return &RedisCaptchaCache{client: client}
}

func buildKey(way domain.VerifyWay, id int64) (string, error) {
	keyPre := utils.GetRedisKey(way.GetKey())
	return fmt.Sprintf("%s:%d", keyPre, id), nil
}

func (r *RedisCaptchaCache) Save(way domain.VerifyWay, value string) (int64, error) {
	id, err := uid.Gen()
	if err != nil {
		return 0, errors.WithStack(err)
	}

	key, err := buildKey(way, id)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	if err := r.client.Set(context.Background(), key, value, way.GetExpire()).Err(); err != nil {
		return 0, errors.WithStack(err)
	}

	return id, nil
}

func (r *RedisCaptchaCache) Verify(way domain.VerifyWay, id int64, value string) error {
	key, err := buildKey(way, id)
	if err != nil {
		return errors.WithStack(err)
	}

	result, err := r.client.Get(context.Background(), key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return codes.ErrCaptchaNotFound
		}
		return err
	}
	if result != value {
		return codes.ErrCaptchaVerifyFailed
	}

	// 验证成功之后删除
	if err := r.Delete(key); err != nil {
		zap.L().Error("删除验证码id缓存失败", zap.String("key", key))
	}

	return nil
}

func (r *RedisCaptchaCache) Delete(key string) error {
	return r.client.Del(context.Background(), key).Err()
}
