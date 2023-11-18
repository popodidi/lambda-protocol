package middlewares

import (
	"errors"
	"net"
	"net/http"
	"os"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// PanicCatcher defines a panic catcher handler.
func PanicCatcher(ctx *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			logger := GetLogger(ctx)

			var brokenPipe bool
			if ne, ok := err.(*net.OpError); ok {
				var se *os.SyscallError
				if errors.As(ne, &se) {
					seStr := strings.ToLower(se.Error())
					if strings.Contains(seStr, "broken pipe") ||
						strings.Contains(seStr, "connection reset by peer") {
						brokenPipe = true
					}
				}
			}

			logger.With(zap.Any("error", err), zap.String("stack", string(debug.Stack()))).Error("panic")
			if brokenPipe {
				// If the connection is dead, we can't write a status to it.
				ctx.Error(err.(error)) //nolint: errcheck
				ctx.Abort()
			} else {
				ctx.AbortWithStatus(http.StatusInternalServerError)
			}
		}
	}()
	// Process Request Chain
	ctx.Next()
}
