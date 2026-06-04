package files

import (
	"bytes"
	"testing"
	"unicode/utf8"
)

var goldens = []struct {
	in  string
	out []byte
}{
	{"A", []byte{0x14}},
	{"{", []byte{194, 183}},
	{"é", []byte{194, 158}},
	{"А", []byte{228, 128, 129}},
	{"я", []byte{228, 131, 180}},
	{"ё", []byte{228, 128, 149}},
	{"Ё", []byte{228, 128, 144}},
}

func TestEncodeDecodeGolden(t *testing.T) {
	for _, g := range goldens {
		got := encodeDecode([]byte(g.in))
		if !bytes.Equal(got, g.out) {
			t.Errorf("encode %q = %v, want %v", g.in, got, g.out)
		}
		if back := encodeDecode(g.out); !bytes.Equal(back, []byte(g.in)) {
			t.Errorf("decode %v = %v, want %q", g.out, back, g.in)
		}
	}
}

func TestEncodeDecodeRoundTrip(t *testing.T) {
	for cp := rune(0x20); cp <= 0x10FFFF; cp++ {
		if cp >= 0xD800 && cp <= 0xDFFF {
			continue
		}
		in := []byte(string(cp))
		if back := encodeDecode(encodeDecode(in)); !bytes.Equal(in, back) {
			t.Fatalf("round-trip failed for U+%04X: got %v", cp, back)
		}
	}
}

func TestEncodeDecodeMixedString(t *testing.T) {
	in := []byte(`{"name":"Гвардеец","note":"café — déjà ✓","n":42}`)
	enc := encodeDecode(in)
	if bytes.Equal(enc, in) {
		t.Fatal("encoding produced identical bytes")
	}
	if !utf8.Valid(enc) {
		t.Fatal("encoded output is not valid UTF-8")
	}
	if back := encodeDecode(enc); !bytes.Equal(in, back) {
		t.Fatalf("round-trip failed: %q", back)
	}
}
