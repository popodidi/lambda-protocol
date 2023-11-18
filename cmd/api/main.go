package main

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/popodidi/lambda-protocol/internal/app/api"
	"github.com/popodidi/lambda-protocol/pkg/app"
	"github.com/popodidi/lambda-protocol/pkg/log"
	"github.com/popodidi/lambda-protocol/pkg/version"
)

func main() {
	sync, err := log.Init(log.Config{
		Name:   "lambda-protocol.api",
		Level:  zapcore.InfoLevel,
		Stdout: true,
		File:   "log/lambda-protocol/api.log",
	})
	if err != nil {
		panic(err)
	}
	defer sync()
	logger := zap.L()
	logger.Info("Hello, world!", zap.String("version", version.FullVersion()))

	ctx := app.GraceCtx(context.Background())
	app := api.New(api.Config{
		Port: 9999,
	})
	err = app.Start(ctx)
	if err != nil {
		logger.Fatal("app.Start", zap.Error(err))
	}
}
