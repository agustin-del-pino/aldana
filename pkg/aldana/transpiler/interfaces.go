package transpiler

// Transpiler provides a convertor from AST to bytes.
// Where T is the type for the AST Nodes.
type Transpiler[T any] interface {
	// Transpile returns the bytes of the transpiled node, or nil and an error.
	Transpile(n T) ([]byte, error)
}
