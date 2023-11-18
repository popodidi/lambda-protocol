package lambda

import (
	"context"

	domain "github.com/popodidi/lambda-protocol/internal/domain/lambda"
)

type IPFSRepository struct {
}

func NewIPFSRepository() domain.Repository {
	return &IPFSRepository{}
}

func (r *IPFSRepository) GetLambada(ctx context.Context, hash string) (*Lambda, error) {
}
