package lambda

import (
	"context"

	domain "github.com/popodidi/lambda-protocol/internal/domain/lambda"
)

func NewRuntimeV8() domain.Runtime {
	return &runtimeV8{}
}

type runtimeV8 struct{}

func (r *runtimeV8) Exec(ctx context.Context, code string, input []byte) ([]byte, error) {
	return nil, nil
}
