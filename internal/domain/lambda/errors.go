package lambda

import "errors"

var (
	ErrInvalidLambda  = errors.New("invalid lambda")
	ErrLambdaNotFound = errors.New("lambda not found")

	ErrRuntimeNotFound = errors.New("runtime not found")
)
