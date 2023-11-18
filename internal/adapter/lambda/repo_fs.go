package lambda

import (
	"context"

	domain "github.com/popodidi/lambda-protocol/internal/domain/lambda"
)

type FSRepository struct {
}

func NewFSRepository() domain.Repository {
	return &FSRepository{}
}

func (r *FSRepository) GetLambdaByName(ctx context.Context, name string) (*domain.Lambda, error) {
	return nil, nil
}
