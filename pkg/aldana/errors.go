package aldana

import (
	"errors"
	"fmt"

	"github.com/agustin-del-pino/aldana/pkg/aldana/lexer"
)

var (
	ErrNoTokenToParser        = errors.New("no token were given to parse")
	ErrNotFoundRootParserRule = errors.New("the root parser rule was not found")
	ErrNotFundParserRule      = errors.New("the parser rule was not found")
	ErrEmptyBytes             = errors.New("the no bytes resulted after the transpilation")
)

func GetLexerError(err error, c lexer.Cursor) error {
	ln, cl := c.GetPosition()
	return fmt.Errorf("%s %s at: line %d column %d", err, string(c.GetChar()), ln, cl)
}
