package lexer

import "errors"

var (
	// ErrUnexpectedChar is returned when a char is not expected for any lexical rule.
	ErrUnexpectedChar = errors.New("unexpected char, cannot create or include in a token")
)
