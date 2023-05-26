package aldana

import "github.com/agustin-del-pino/aldana/pkg/aldana/lexer"

// defaultCursor implements lexer.Cursor.
type defaultCursor struct {
	content []byte
	length  int
	column  int
	line    int
	char    byte
}

func (c *defaultCursor) HasChar() bool {
	return c.column < c.length
}

func (c *defaultCursor) GetChar() byte {
	return c.char
}

func (c *defaultCursor) GetPosition() (int, int) {
	return c.column, c.line
}

func (c *defaultCursor) Next() {
	if c.column < c.length {
		c.char = c.content[c.column]
		c.column += 1
	}
}

func (c *defaultCursor) AddLine(l int) {
	c.line += l
}

// NewCursor returns the default implementation of lexer.Cursor.
//
// # About this implementation
//   - the cursor can read when the column value is less than the length of the reading bytes.
//
// # Example
//
//	cur := NewCursor([]byte("123 456 799"))
//
//	b, _ := os.ReadFile("text.txt")
//	curF := NewCursor(b)
func NewCursor(c []byte) lexer.Cursor {
	return &defaultCursor{
		content: c,
		length:  len(c),
	}
}
