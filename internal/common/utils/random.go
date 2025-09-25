package utils

import (
	"math/rand"
	"time"
)

func GenRandomCodeForJWT(length int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := make([]byte, length)
	for i := 0; i < length; i++ {
		code[i] = '0' + byte(r.Intn(10))
	}
	return string(code)
}
