package aldana

import (
	"testing"

	"github.com/agustin-del-pino/aldana/pkg/aldana/lexer"
	"github.com/stretchr/testify/assert"
)

/*
Given: a n-len byte array.
When: advances the column and gets the char
Then: returns the current char.
*/
func TestCursor_GetChar_with_no_reaching(t *testing.T) {
	t.Run("no next", func(t *testing.T) {
		// arrange
		cur := setUpCursor([]byte("123456789"))

		// act
		ch := cur.GetChar()

		// assert
		assert.Equal(t, byte(0x00), ch)
	})
	t.Run("one next", func(t *testing.T) {
		// arrange
		cur := setUpCursor([]byte("123456789"))

		// act
		cur.Next()
		ch := cur.GetChar()

		// assert
		assert.Equal(t, byte('1'), ch)
	})
	t.Run("next to penultimate", func(t *testing.T) {
		// arrange
		cur := setUpCursor([]byte("123456789"))

		// act
		for i := 0; i < 8; i++ {
			cur.Next()
		}

		ch := cur.GetChar()

		// assert
		assert.Equal(t, byte('8'), ch)
	})

	t.Run("next to last one", func(t *testing.T) {
		// arrange
		cur := setUpCursor([]byte("123456789"))

		// act
		for i := 0; i < 10; i++ {
			cur.Next()
		}

		ch := cur.GetChar()

		// assert
		assert.Equal(t, byte('9'), ch)
	})

	t.Run("over-reaching", func(t *testing.T) {
		// arrange
		cur := setUpCursor([]byte("123456789"))

		// act
		for i := 0; i < 50; i++ {
			cur.Next()
		}

		ch := cur.GetChar()

		// assert
		assert.Equal(t, byte('9'), ch)
	})
}

/*
Given: a n-len byte array.
When: advances n < n-len column and performs the HasChar
Then: returns true.
*/
func TestCursor_HasToken_with_no_reaching(t *testing.T) {
	t.Run("no next", func(t *testing.T) {
		// arrange
		cur := setUpCursor([]byte("123456789"))

		// act
		hc := cur.HasChar()

		// assert
		assert.True(t, hc)
	})
	t.Run("one next", func(t *testing.T) {
		// arrange
		cur := setUpCursor([]byte("123456789"))

		// act
		cur.Next()
		hc := cur.HasChar()

		// assert
		assert.True(t, hc)
	})
	t.Run("next to penultimate", func(t *testing.T) {
		// arrange
		cur := setUpCursor([]byte("123456789"))

		// act
		for i := 0; i < 8; i++ {
			cur.Next()
		}
		hc := cur.HasChar()

		// assert
		assert.True(t, hc)
	})
}

/*
Given: a n-len byte array.
When: reaches n column and performs the HasChar.
Then: returns false.
*/
func TestCursor_HasToken_with_reaching(t *testing.T) {
	t.Run("0-len", func(t *testing.T) {
		// arrange
		cur := setUpCursor([]byte(""))

		// act
		hc := cur.HasChar()

		// assert
		assert.False(t, hc)
	})

	t.Run("1-len", func(t *testing.T) {
		// arrange
		cur := setUpCursor([]byte("a"))

		// act
		cur.Next()
		hc := cur.HasChar()

		// assert
		assert.False(t, hc)
	})

	t.Run("large-len", func(t *testing.T) {
		// arrange
		cur := setUpCursor([]byte("qwertyuiopasdfghjklñzxcvbnm,"))

		// act
		for i := 0; i < 29; i++ {
			cur.Next()
		}

		hc := cur.HasChar()

		// assert
		assert.False(t, hc)
	})

	t.Run("large-len and overreach", func(t *testing.T) {
		// arrange
		cur := setUpCursor([]byte("qwertyuiopasdfghjklñzxcvbnm,"))

		// act
		for i := 0; i < 50; i++{
			cur.Next()
		}

		hc := cur.HasChar()

		// assert
		assert.False(t, hc)
	})
}

func setUpCursor(c []byte) lexer.Cursor {
	return NewCursor(c)
}
