package api

import (
	"context"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	svc "github.com/popodidi/lambda-protocol/internal/service/lambda"
	"github.com/popodidi/lambda-protocol/pkg/api/middlewares"
	"github.com/popodidi/lambda-protocol/pkg/api/protocol"
)

const (
	paramLambdaName = "lambda_name"
)

type Config struct {
	Port    int
	Service svc.Service
}

func New(config Config) *App {
	return &App{config}
}

type App struct {
	Config
}

func (a *App) Start(ctx context.Context) error {
	// start a gin HTTP server
	// Create our HTTP Router
	router := gin.New()

	// Configure HTTP Router Settings
	router.RedirectTrailingSlash = true
	router.RedirectFixedPath = false
	router.HandleMethodNotAllowed = false
	router.ForwardedByClientIP = true
	router.AppEngine = false
	router.UseRawPath = false
	router.UnescapePathValues = true
	router.ContextWithFallback = true

	router.Use(middlewares.PanicCatcher)
	router.Use(cors.Default())
	router.Use(middlewares.CtxLogger)

	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	router.POST(fmt.Sprintf("lambda/:%s", paramLambdaName), a.exec)

	// Setup HTTP Server
	server := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", a.Port),
		Handler: router,
	}

	// Start Running HTTP Server.
	go server.ListenAndServe()
	<-ctx.Done()
	return server.Shutdown(ctx)
}

func (a *App) exec(ctx *gin.Context) {
	logger := middlewares.GetLogger(ctx)

	// parse inputs
	lambdaName := ctx.Param(paramLambdaName)
	var (
		req        ExecRequest
		inputBytes []byte
	)

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		logger.Error("c.ShouldBindJSON", zap.Error(err))
		protocol.CtxErr(ctx, http.StatusBadRequest, "")
		return
	}
	if len(req.Input) > 0 {
		inputBytes, err = hex.DecodeString(req.Input)
		if err != nil {
			logger.Error("hex.DecodeString", zap.Error(err))
			protocol.CtxErr(ctx, http.StatusBadRequest, "")
			return
		}
	}

	// execute lambda
	var res ExecResponse
	output, err := a.Service.Execute(ctx, lambdaName, inputBytes)
	if len(output) > 0 {
		res.Output = hex.EncodeToString(output)
	}
	if err != nil {
		res.Error = err.Error()
	}
	protocol.CtxOKResult(ctx, res)
}
