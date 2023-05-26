package aldana

import "github.com/agustin-del-pino/aldana/pkg/aldana/transpiler"

type NodePredicator[T any] func(n T) bool

type ByteTranspiler[T any] func(n T, t transpiler.Transpiler[T]) ([]byte, error)

type TranspileRule[T any] func() (NodePredicator[T], ByteTranspiler[T])

type TranspilerOptions[T any] struct {
	TranspileRules []TranspileRule[T]
}

type defaultTranspiler[T any] struct {
	ops *TranspilerOptions[T]
}

// Transpile returns the bytes of the transpiled node, or nil and an error.
func (t *defaultTranspiler[T]) Transpile(n T) ([]byte, error) {
	var b []byte

	for _, r := range t.ops.TranspileRules {
		np, bt := r()
		if !np(n) {
			continue
		}
		by, bErr := bt(n, t)

		if bErr != nil {
			return nil, bErr
		}

		b = append(b, by...)
	}

	if b != nil {
		if len(b) == 0 {
			return nil, ErrEmptyBytes
		}
	} else {
		return nil, ErrEmptyBytes
	}

	return b, nil
}

// NewTranspiler returns the default implementation for transpiler.Transpiler.
// Where T is the type for the nodes.
//
// # Example
//
//	func transpileSum(n *Node, t transpiler.Transpiler[*Node]) ([]byte, error) {
//		var s []byte
//		for i :=0; i < len(n.Children)-1; i++ {
//			b, err := t.Transpile(n.Children[i])
//			if err != nil {
//				return nil, err
//			}
//			s = append(s, b...)
//			s = append(s, '+')
//		}
//
//		b, err := t.Transpile(n.LastChildren)
//		if err != nil {
//			return nil, err
//		}
//		s = append(s, b...)
//		return s, nil
//	}
//
//	func transpileSum(n *Node, t transpiler.Transpiler[*Node]) ([]byte, error) {
//		return n.Token.Value, nil
//	}
//
//	trs := NewTranspiler[*Node](&TranspilerOptions{
//		TranspileRules: []TranspileRules{
//			NewTranspileRule(IsSumNode, transpileSum),
//			NewTranspileRule(IsNumNode, transpileNum),
//		}
//	})
func NewTranspiler[T any](ops *TranspilerOptions[T]) transpiler.Transpiler[T] {
	return &defaultTranspiler[T]{
		ops: ops,
	}
}

func NewTranspileRule[T any](p NodePredicator[T], b ByteTranspiler[T]) TranspileRule[T] {
	return func() (NodePredicator[T], ByteTranspiler[T]) {
		return p, b
	}
}
