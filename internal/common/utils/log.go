package utils

import (
	"fmt"
)

func Println(n int) {
	for _ = range n {
		fmt.Println()
	}
}
