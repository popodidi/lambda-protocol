package lambda

import (
	"context"
	"encoding/json"
	"fmt"

	"go.uber.org/zap"
	v8 "rogchap.com/v8go"

	domain "github.com/popodidi/lambda-protocol/internal/domain/lambda"
	"github.com/popodidi/lambda-protocol/pkg/log"
)

func NewRuntimeV8() domain.Runtime {
	return &runtimeV8{}
}

type runtimeV8 struct{}

func (r *runtimeV8) Exec(ctx context.Context, code string, input []byte) ([]byte, error) {
	// prepare logger
	var logger *zap.Logger
	ctx, logger = log.Context(ctx)

	// validate if input is a valid JSON
	var inputStr string
	var err error
	if len(input) > 0 {
		err = json.Unmarshal(input, &struct{}{})
		if err != nil {
			logger.Error("json.Unmarshal", zap.Error(err), zap.String("input", string(input)))
			return nil, err
		}
		inputStr = string(input)
	}

	v8ctx := v8.NewContext()
	_, err = v8ctx.RunScript(code, "lambda.js")
	if err != nil {
		logger.Error("v8ctx.RunScript", zap.Error(err), zap.String("code", code))
		return nil, err
	}
	value, err := v8ctx.RunScript(fmt.Sprintf("main(%s)", inputStr), "lambda.js")
	if err != nil {
		logger.Error("v8ctx.RunScript", zap.Error(err), zap.String("input", inputStr))
		return nil, err
	}
	return value.MarshalJSON()
}
