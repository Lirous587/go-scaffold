package utils

const (
	Prefix = "scaffold:" //项目key前缀
)

func GetRedisKey(key string) string {
	return Prefix + key
}
