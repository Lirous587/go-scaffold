package utils

const (
	Prefix = "your-prefix" //项目key前缀
)

func GetRedisKey(key string) string {
	// panic("请先设置项目前缀，然后移除此行代码")
	return Prefix + key
}
