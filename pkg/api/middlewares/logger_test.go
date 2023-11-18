package middlewares

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/popodidi/lambda-protocol/pkg/log"
)

func TestMiddlewareFirst(t *testing.T) {
	log.Init(log.Config{})
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = &http.Request{
		URL: &url.URL{},
	}

	logger := GetLogger(ctx)
	require.NotNil(t, logger)

	_, logLogger := log.Context(ctx.Request.Context(), "test")
	require.NotNil(t, logLogger)

	logLoggerEntry := logLogger.Check(zapcore.InfoLevel, "")
	loggerEntry := logger.Check(zapcore.InfoLevel, "")

	require.Equal(t, logLoggerEntry.Entry.LoggerName, loggerEntry.Entry.LoggerName)
}

func TestLoggerFirst(t *testing.T) {
	log.Init(log.Config{})
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = &http.Request{
		URL: &url.URL{},
	}

	reqCtx, logLogger := log.Context(ctx.Request.Context(), "test")
	ctx.Request = ctx.Request.WithContext(reqCtx)
	require.NotNil(t, logLogger)

	logger := GetLogger(ctx)
	require.NotNil(t, logger)

	logLoggerEntry := logLogger.Check(zapcore.InfoLevel, "")
	loggerEntry := logger.Check(zapcore.InfoLevel, "")

	require.Equal(t, logLoggerEntry.Entry.LoggerName, loggerEntry.Entry.LoggerName)
}
