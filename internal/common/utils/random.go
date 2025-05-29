package utils

import (
	"math/rand"
	"time"
)

const jwtRandomCodeLength = 20

func GenRandomCodeForJWT() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := make([]byte, jwtRandomCodeLength)
	for i := 0; i < jwtRandomCodeLength; i++ {
		code[i] = '0' + byte(r.Intn(10))
	}
	return string(code)
}
