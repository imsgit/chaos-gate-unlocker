package files

import (
	"bytes"
	"unicode/utf8"
)

func encodeDecode(data []byte) []byte {
	var result bytes.Buffer
	result.Grow(len(data))

	for len(data) > 0 {
		r, size := utf8.DecodeRune(data)
		data = data[size:]

		if s := swapNibbles(r); utf8.ValidRune(s) {
			r = s
		}
		result.WriteRune(r)
	}

	return result.Bytes()
}

func swapNibbles(r rune) rune {
	v := uint32(r)
	return rune((v&0x0F0F0F0F)<<4 | (v&0xF0F0F0F0)>>4)
}
