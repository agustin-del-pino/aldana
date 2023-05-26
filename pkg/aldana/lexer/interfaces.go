package lexer

// Cursor provides a bytes reader.
type Cursor interface {
	// HasChar returns a boolean that indicates whether still bytes to read.
	HasChar() bool
	// GetChar returns the current char where the cursor is positioned.
	GetChar() byte
	// GetPosition returns the column and line where the cursor is positioned.
	GetPosition() (int, int)
	// Next advances the cursor to next column.
	Next()
	// AddLine adds a line. Where l is the number of line to add.
	AddLine(l int)
}

// Lexer provides a lexer-processor that tokenize an input into an array of tokens.
// Where T is the type of the tokens.
//
// # Recommendations:
//
//   - In case of use some kind of struct, T must be a pointer: Lexer[*MyTokenStruct]. Avoid Lexer[MyTokenStruct].
//   - In case of use some kind of interface, T must not be a pointer: Lexer[IToken]. Avoid Lexer[*IToken].
type Lexer[T any] interface {
	// Tokenize returns a slice of tokens, or nil and a lexer-error.
	Tokenize(c Cursor) ([]T, error)
}
