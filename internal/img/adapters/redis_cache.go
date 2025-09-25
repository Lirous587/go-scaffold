package adapters

import (
	"scaffold/internal/common/utils"
	"scaffold/internal/img/domain"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"os"
	"strconv"
	"strings"
	"time"
)

type RedisCache struct {
	client *redis.Client
}

func NewImgRedisCache() domain.ImgMsgQueue {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	password := os.Getenv("REDIS_PASSWORD")
	dbStr := os.Getenv("REDIS_DB")
	poolSizeStr := os.Getenv("REDIS_POOL_SIZE")

	db, _ := strconv.Atoi(dbStr)
	poolSize, _ := strconv.Atoi(poolSizeStr)

	addr := host + ":" + port

	client := redis.NewClient(&redis.Options{
		Addr:		addr,
		DB:		db,
		Password:	password,
		PoolSize:	poolSize,
	})

	// 可选：ping 检查连接
	if err := client.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}

	return &RedisCache{client: client}
}

const (
	keyImgDeleteQueueKey	= "img:delete"
	deleteImgExpire		= time.Hour * 24 * 7
	//deleteImgExpire = time.Second * 20
)

// AddToDeleteQueue 软删除时写入 redis 并设置过期
func (c *RedisCache) AddToDeleteQueue(imgID int64) error {
	key := fmt.Sprintf("%s:%d", utils.GetRedisKey(keyImgDeleteQueueKey), imgID)
	return c.client.SetEx(context.Background(), key, "", deleteImgExpire).Err()
}

// ListenDeleteQueue 后台监听 key 过期事件
func (c *RedisCache) ListenDeleteQueue(onExpire func(imgID int64)) {
	pubsub := c.client.PSubscribe(context.Background(), "__keyevent@0__:expired")
	defer pubsub.Close()	// 确保资源释放

	preKey := utils.GetRedisKey(keyImgDeleteQueueKey) + ":"

	for msg := range pubsub.Channel() {
		// 解析 key
		if strings.HasPrefix(msg.Payload, preKey) {
			idStr := strings.TrimPrefix(msg.Payload, preKey)
			if imgID, err := strconv.ParseInt(idStr, 10, 64); err == nil {
				onExpire(imgID)
			}
		}
	}
}

// SendDeleteMsg 发送删除消息
func (c *RedisCache) SendDeleteMsg(imgID int64) error {
	key := fmt.Sprintf("%s:%d", utils.GetRedisKey(keyImgDeleteQueueKey), imgID)

	// 使用 SETEX 设置过期时间为 1 毫秒，几乎立即过期
	return c.client.SetEx(context.Background(), key, "manual_delete", time.Millisecond).Err()
}

// RemoveFromDeleteQueue 从删除队列中移除指定图片ID
func (c *RedisCache) RemoveFromDeleteQueue(imgID int64) error {
	key := fmt.Sprintf("%s:%d", utils.GetRedisKey(keyImgDeleteQueueKey), imgID)
	return c.client.Del(context.Background(), key).Err()
}
