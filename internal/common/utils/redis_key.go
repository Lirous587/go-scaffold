package utils

const (
	Prefix = "your-prefix" //项目key前缀
)

func GetRedisKey(key string) string {
	// panic("请确保设置项目前缀")
	return Prefix + key
}
