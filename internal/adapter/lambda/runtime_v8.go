package lambda

import (
	"context"

	v8 "rogchap.com/v8go"

	domain "github.com/popodidi/lambda-protocol/internal/domain/lambda"
)

func NewRuntimeV8() domain.Runtime {
	return &runtimeV8{}
}

type runtimeV8 struct{}

func (r *runtimeV8) Exec(ctx context.Context, code string, input []byte) ([]byte, error) {
	v8ctx := v8.NewContext()
	value, err := v8ctx.RunScript(code, "lambda.js")
	if err != nil {
		return nil, err
	}
	return value.MarshalJSON()
}
