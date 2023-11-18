package protocol

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetPathInt(ctx *gin.Context, pathParam string) (int, error) {
	paramStr := ctx.Param(pathParam)
	if len(paramStr) == 0 {
		return 0, errors.New("empty_param_id")
	}
	param, err := strconv.Atoi(paramStr)
	if err != nil {
		return 0, err
	}
	return param, nil
}
