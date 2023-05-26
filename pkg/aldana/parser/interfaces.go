package parser

// Reader provides a token inspector.
// Where T is the type for the Tokens.
//
// # Recommendations
//   - In case of use some kind of struct, T must be a pointer: Reader[*MyTokenStruct]. Avoid Reader[MyTokenStruct].
//   - In case of use some kind of interface, T must not be a pointer: Reader[IToken]. Avoid Reader[*IToken].
type Reader[T any] interface {
	// HasTokens returns a boolean that indicates whether still tokens to read.
	HasTokens() bool
	// GetToken returns the current token where the reader is positioned.
	GetToken() T
	// Next advances the reader to next token.
	Next()
}

// Parser provides an analyzer of tokens.
// Where Tt is the type for the Tokens, adn Tn is the type for the Node.
//
// # Recommendations
//   - In case of use some kind of struct, Tt or Tn must be a pointer: Parser[*MyTokenStruct, *MyNode]. Avoid Parser[MyTokenStruct, Node].
//   - In case of use some kind of interface, Tt or Tn must not be a pointer: Parser[IToken, INode]. Avoid Parser[*IToken, *INode].
type Parser[Tt any, Tn any] interface {
	// Parse returns the result node from the analyze of the tokens, or nil and an error.
	Parse(r Reader[Tt]) (Tn, error)
}
