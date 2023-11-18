package main

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	adapter "github.com/popodidi/lambda-protocol/internal/adapter/lambda"
	"github.com/popodidi/lambda-protocol/internal/app/api"
	domain "github.com/popodidi/lambda-protocol/internal/domain/lambda"
	service "github.com/popodidi/lambda-protocol/internal/service/lambda"
	"github.com/popodidi/lambda-protocol/pkg/app"
	"github.com/popodidi/lambda-protocol/pkg/log"
	"github.com/popodidi/lambda-protocol/pkg/version"
)

func main() {
	// init logger
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

	// prepare context
	ctx := app.GraceCtx(context.Background())

	// prepare service
	runtimeV8 := adapter.NewRuntimeV8()
	repo := adapter.NewFSRepository()
	svc := service.NewService(
		repo,
		map[domain.RuntimeType]domain.Runtime{
			domain.RuntimeTypeV8: runtimeV8,
		},
	)

	app := api.New(api.Config{
		Port:    9999,
		Service: svc,
	})
	err = app.Start(ctx)
	if err != nil {
		logger.Fatal("app.Start", zap.Error(err))
	}
}
