package ranges

// ByteRange is an alias for a function that take a byte and indicates with a boolean whether the byte is in-range or not.
type ByteRange func(b byte) bool

// ByteSet returns a ByteRange that indicates whether the byte is included in the set or not.
func ByteSet(bs ...byte) ByteRange {
	return func(b byte) bool {
		for _, v := range bs {
			if b == v {
				return true
			}
		}
		return false
	}
}

// ByteBounded returns a ByteRange that indicates whether the byte is bounded in the range or not.
func ByteBounded(sb byte, eb byte) ByteRange {
	return func(b byte) bool {
		return b >= sb && b <= eb
	}
}

// ByteSingle returns a ByteRange that indicates whether the byte is equal to bs or not.
func ByteSingle(bs byte) ByteRange {
	return func(b byte) bool {
		return b == bs
	}
}

func RangeByteOfRange(br ...ByteRange) ByteRange {
	return func(b byte) bool {
		for _, r := range br {
			if r(b) {
				return true
			}
		}
		return false
	}
}
