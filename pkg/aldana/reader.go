package aldana

import "github.com/agustin-del-pino/aldana/pkg/aldana/parser"

type defaultReader[T any] struct {
	tokens   []T
	token    T
	length   int
	position int
}

func (r *defaultReader[T]) HasTokens() bool {
	return r.position < r.length
}

func (r *defaultReader[T]) GetToken() T {
	return r.token
}

func (r *defaultReader[T]) Next() {
	if r.position < r.length {
		r.token = r.tokens[r.position]
		r.position += 1
	}
}

func NewReader[T any](t []T) parser.Reader[T] {
	return &defaultReader[T]{
		tokens: t,
		length: len(t),
	}
}
