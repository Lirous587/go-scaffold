package utils

const (
	Prefix = "blog-v4:" //项目key前缀
)

func GetRedisKey(key string) string {
	return Prefix + key
}
