package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/popodidi/lambda-protocol/pkg/log"
	"go.uber.org/zap"
)

func GraceCtx(parent context.Context) (ctx context.Context) {
	ctx, cancel := context.WithCancel(parent)
	go func() {
		quitChan := make(chan os.Signal, 1)
		signal.Notify(quitChan, os.Interrupt, syscall.SIGTERM)
		sig := <-quitChan
		_, logger := log.Context(parent, "app")
		logger.With(zap.String("signal", sig.String())).Warn("graceful_shutdown")
		cancel()
	}()
	return ctx
}
