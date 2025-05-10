package utils

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetID(ctx *gin.Context) (uint, error) {
	idStr := ctx.Param("id")
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, err
	}
	id := uint(idInt)
	return id, nil
}
