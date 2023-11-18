package lambda

import "context"

type Lambda struct {
	Name     string
	Metadata Metadata
	Code     string
}

type Metadata struct {
	RuntimeType RuntimeType `json:"runtime_type"`
	CodeFile    string      `json:"code_file"`
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
