package aldana

import "errors"

var (
	ErrNoTokenToParser        = errors.New("no token were given to parse")
	ErrNotFoundRootParserRule = errors.New("the root parser rule was not found")
	ErrNotFundParserRule      = errors.New("the parser rule was not found")
	ErrEmptyBytes             = errors.New("the no bytes resulted after the transpilation")
)
