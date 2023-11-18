package lambda

import (
	"context"
	"fmt"

	domain "github.com/popodidi/lambda-protocol/internal/domain/lambda"
)

type Service interface {
	Execute(ctx context.Context, lambdaName string, input []byte) ([]byte, error)
}

type service struct {
	repo     domain.Repository
	runtimes map[domain.RuntimeType]domain.Runtime
}

func NewService() Service {
	return &service{}
}

func (s *service) Execute(ctx context.Context, lambdaName string, input []byte) ([]byte, error) {
	lambda, err := s.repo.GetLambdaByName(ctx, lambdaName)
	if err != nil {
		return nil, err
	}
	runtime, ok := s.runtimes[lambda.Metadata.RuntimeType]
	if !ok {
		return nil, fmt.Errorf(
			"runtime %s not found; %w",
			lambda.Metadata.RuntimeType, domain.ErrRuntimeNotFound,
		)
	}
	return runtime.Exec(ctx, lambda.Code, input)
}
