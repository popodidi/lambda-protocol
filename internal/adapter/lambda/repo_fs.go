package lambda

import (
	"context"
	"encoding/json"
	"os"

	"go.uber.org/zap"

	domain "github.com/popodidi/lambda-protocol/internal/domain/lambda"
	"github.com/popodidi/lambda-protocol/pkg/log"
)

type FSRepository struct {
	root string
}

func NewFSRepository(root string) domain.Repository {
	return &FSRepository{
		root: root,
	}
}

func (r *FSRepository) GetLambdaByName(ctx context.Context, name string) (*domain.Lambda, error) {
	// prepare logger
	var logger *zap.Logger
	ctx, logger = log.Context(ctx)

	// read metadata.json
	metaBytes, err := os.ReadFile(r.root + "/" + name + "/metadata.json")
	if err != nil {
		logger.Error("os.ReadFile", zap.Error(err))
		return nil, err
	}

	// unmarshall metadata.json
	var meta domain.Metadata
	err = json.Unmarshal(metaBytes, &meta)
	if err != nil {
		logger.Error("json.Unmarshal", zap.Error(err))
		return nil, err
	}

	// read code
	codeFile := r.root + "/" + name + "/" + meta.CodeFile
	codeBytes, err := os.ReadFile(codeFile)
	if err != nil {
		logger.Error("os.ReadFile", zap.Error(err))
		return nil, err
	}

	return &domain.Lambda{
		Name:     name,
		Metadata: meta,
		Code:     string(codeBytes),
	}, nil
}
