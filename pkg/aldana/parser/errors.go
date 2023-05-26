package parser

import "errors"

var (
	ErrInvalidSyntax  = errors.New("the syntax is invalid")
	ErrUnhandledToken = errors.New("unhandled token")
)
