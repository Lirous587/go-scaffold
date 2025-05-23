package utils

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetID(ctx *gin.Context) (int, error) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, err
	}
	return id, nil
}
