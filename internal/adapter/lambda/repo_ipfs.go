package lambda

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"sync"

	"github.com/ipfs/boxo/files"
	"github.com/ipfs/boxo/path"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/kubo/client/rpc"
	"go.uber.org/zap"

	domain "github.com/popodidi/lambda-protocol/internal/domain/lambda"
	"github.com/popodidi/lambda-protocol/pkg/log"
)

type IPFSRepository struct {
	ipfsAPI string
	cid     string

	node struct {
		sync.Once
		*rpc.HttpApi
	}
}

func NewIPFSRepository(ipfsAPI, cid string) domain.Repository {
	return &IPFSRepository{
		ipfsAPI: ipfsAPI,
		cid:     cid,
	}
}

func (r *IPFSRepository) getClient() (*rpc.HttpApi, error) {
	var err error
	r.node.Do(func() {
		r.node.HttpApi, err = rpc.NewURLApiWithClient(r.ipfsAPI, &http.Client{})
	})
	if err != nil {
		// reset node
		r.node.HttpApi = nil
		r.node.Once = sync.Once{}
		return nil, err
	}
	return r.node.HttpApi, nil
}

func (r *IPFSRepository) GetLambdaByName(ctx context.Context, name string) (*domain.Lambda, error) {
	// prepare logger
	var logger *zap.Logger
	ctx, logger = log.Context(ctx)

	// get client
	client, err := r.getClient()
	if err != nil {
		logger.Error("r.getClient", zap.Error(err))
		return nil, err
	}
	fs := client.Unixfs()
	lambdaPath, err := path.NewPath("/ipfs/" + r.cid)
	if err != nil {
		logger.Error("path.NewPath", zap.Error(err))
		return nil, err
	}

	// build lambda map
	lambdas := make(map[string]cid.Cid)
	entries, err := fs.Ls(ctx, lambdaPath)
	if err != nil {
		logger.Error("fs.Ls", zap.Error(err))
		return nil, err
	}
	for entry := range entries {
		lambdas[entry.Name] = entry.Cid
	}

	// find lambda
	lambdaCid, ok := lambdas[name]
	if !ok {
		logger.Error("lambda not found", zap.String("name", name))
		return nil, domain.ErrLambdaNotFound
	}

	// build lambda file map
	lambdaFiles := make(map[string]cid.Cid)
	entries, err = fs.Ls(ctx, path.FromCid(lambdaCid))
	if err != nil {
		logger.Error("fs.Ls", zap.Error(err))
		return nil, err
	}
	for entry := range entries {
		lambdaFiles[entry.Name] = entry.Cid
	}

	// read metadata.json
	metaCid, ok := lambdaFiles["metadata.json"]
	if !ok {
		logger.Error("metadata.json not found")
		return nil, domain.ErrInvalidLambda
	}
	file, err := fs.Get(ctx, path.FromCid(metaCid))
	if err != nil {
		logger.Error("fs.Get", zap.Error(err))
		return nil, err
	}
	metaBytes, err := io.ReadAll(file.(files.File))
	if err != nil {
		logger.Error("io.ReadAll", zap.Error(err))
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
	codeCid, ok := lambdaFiles[meta.CodeFile]
	if !ok {
		logger.Error("code file not found")
		return nil, domain.ErrInvalidLambda
	}
	file, err = fs.Get(ctx, path.FromCid(codeCid))
	if err != nil {
		logger.Error("fs.Get", zap.Error(err))
		return nil, err
	}
	codeBytes, err := io.ReadAll(file.(files.File))
	if err != nil {
		logger.Error("io.ReadAll", zap.Error(err))
		return nil, err
	}

	return &domain.Lambda{
		Name:     name,
		Metadata: meta,
		Code:     string(codeBytes),
	}, nil
}
