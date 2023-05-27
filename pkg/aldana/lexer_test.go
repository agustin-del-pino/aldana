package aldana

import (
	"testing"

	"github.com/agustin-del-pino/aldana/pkg/aldana/lexer"
	"github.com/agustin-del-pino/aldana/pkg/aldana/ranges"
	"github.com/stretchr/testify/assert"
)

type token struct {
	Type  string
	Value []byte
}

/*
Given: a lex-rule and a n-len cursor with all ignorable chars.
When: tokenizes the chars.
Then: returns a token slice of 0-len and no error.
*/
func TestLexer_Tokenize_with_all_ignorable_chars(t *testing.T) {
	// arrange
	lex := setUpLexer(mockLexerRule())

	// act
	tks, err := lex.Tokenize(mockCursor([]byte("    ")))

	// assert
	assert.NoError(t, err)
	assert.Len(t, tks, 0)
}

/*
Given: a lex-rule and a n-len cursor with non acceptable chars.
When: tokenizes the chars.
Then: returns ErrUnexpectedChar a no tokens.
*/
func TestLexer_Tokenize_with_non_acceptable_chars(t *testing.T) {
	t.Run("one non acceptable char", func(t *testing.T) {
		// arrange
		lex := setUpLexer(mockLexerRule())

		// act
		tks, err := lex.Tokenize(mockCursor([]byte("123A4567890")))

		// assert
		assert.ErrorIs(t, err, lexer.ErrUnexpectedChar)
		assert.Nil(t, tks)
	})

	t.Run("many non acceptable chars", func(t *testing.T) {
		// arrange
		lex := setUpLexer(mockLexerRule())

		// act
		tks, err := lex.Tokenize(mockCursor([]byte("1%#$%23 45AAA6 78TERYRTUTYI/**9 0")))

		// assert
		assert.ErrorIs(t, err, lexer.ErrUnexpectedChar)
		assert.Nil(t, tks)
	})
}

/*
Given: a lex-rule and a n-len cursor with only acceptable chars.
When: tokenizes the chars.
Then: returns an slice of tokens and no error.
*/
func TestLexer_Tokenize_with_only_acceptable_chars(t *testing.T) {
	t.Run("one token", func(t *testing.T) {
		// arrange
		lex := setUpLexer(mockLexerRule())

		// act
		tks, err := lex.Tokenize(mockCursor([]byte("1234567890")))

		// assert
		assert.NoError(t, err)
		assert.Len(t, tks, 1)
		assert.Equal(t, "1234567890", string(tks[0].Value))
		assert.Equal(t, "num", string(tks[0].Type))
	})

	t.Run("many token", func(t *testing.T) {
		// arrange
		lex := setUpLexer(mockLexerRule())

		// act
		tks, err := lex.Tokenize(mockCursor([]byte("123 456 789 0")))

		// assert
		assert.NoError(t, err)
		assert.Len(t, tks, 4)
		assert.Equal(t, "123", string(tks[0].Value))
		assert.Equal(t, "456", string(tks[1].Value))
		assert.Equal(t, "789", string(tks[2].Value))
		assert.Equal(t, "0", string(tks[3].Value))

		for i := 0; i < 4; i++ {
			assert.Equal(t, "num", string(tks[i].Type))
		}
	})
}

func mockCursor(b []byte) lexer.Cursor {
	return NewCursor(b)
}

func mockLexerRule() LexicalRule[*token] {
	return func() (ranges.ByteRange, TokenRule[*token]) {
		return ranges.ByteBounded(0x30, 0x39), func(c lexer.Cursor, r ranges.ByteRange) *token {
			t := &token{
				Type: "num",
			}

			for c.HasChar() && r(c.GetChar()) {
				t.Value = append(t.Value, c.GetChar())
				c.Next()
			}

			return t
		}
	}
}

func setUpLexer(r LexicalRule[*token]) lexer.Lexer[*token] {
	return NewLexer(&LexerOptions[*token]{
		Ignore:   IgnoreWhiteSpaces(),
		LexRules: []LexicalRule[*token]{r},
	})
}
