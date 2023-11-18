package lambda

import "context"

type Lambda struct {
	Metadata Metadata
	Code     string
}

type Metadata struct {
	Hash        string
	RuntimeType RuntimeType
}

type Repository interface {
	GetLambdaByName(ctx context.Context, name string) (*Lambda, error)
}

type RuntimeType string

const (
	RuntimeTypeV8 RuntimeType = "v8"
)

type Runtime interface {
	Exec(ctx context.Context, code string, input []byte) ([]byte, error)
}
