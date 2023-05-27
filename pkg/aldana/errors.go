package aldana

import (
	"fmt"

	"github.com/agustin-del-pino/aldana/pkg/aldana/lexer"
)

func GetLexerError(err error, c lexer.Cursor) error {
	ln, cl := c.GetPosition()
	return fmt.Errorf("%s %s at: line %d column %d", err, string(c.GetChar()), ln, cl)
}
