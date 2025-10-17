package utils

import (
	"fmt"
	"os"
	"strconv"
)

func GetEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		panic(fmt.Sprintf("环境变量 %s 未设置或为空", key))
	}
	return val
}

func GetEnvAsInt(key string) int {
	val := os.Getenv(key)
	if val == "" {
		panic(fmt.Sprintf("环境变量 %s 未设置或为空", key))
	}
	intVal, err := strconv.Atoi(val)
	if err != nil {
		panic(fmt.Sprintf("环境变量 %s 无法转换为 int: %v", key, err))
	}
	return intVal
}
