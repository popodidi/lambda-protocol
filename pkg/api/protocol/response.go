package protocol

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool            `json:"success"`
	Error   string          `json:"error,omitempty"`
	Result  json.RawMessage `json:"data,omitempty"`
}

func CtxOK(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, Response{Success: true})
}

func CtxOKResult(ctx *gin.Context, result interface{}) {
	b, err := json.Marshal(result)
	if err != nil {
		CtxErr(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	CtxOKResultRaw(ctx, b)
}

func CtxOKResultRaw(ctx *gin.Context, result []byte) {
	ctx.JSON(http.StatusOK, Response{Success: true, Result: result})
}

func CtxErr(ctx *gin.Context, code int, err string) {
	ctx.JSON(code, Response{Success: false, Error: err})
}
