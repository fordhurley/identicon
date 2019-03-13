package identicon

import (
	"bytes"
	"math/rand"
	"testing"
	"testing/quick"
)

func TestBitReader(t *testing.T) {
	buf := []byte{0xFF, 0, 0x55}
	br := bitReader{r: bytes.NewReader(buf)}
	for i := 0; i < 8; i++ {
		x, err := br.read()
		if err != nil {
			panic(err)
		}
		if !x {
			t.Errorf("expected all bits to be 1")
		}
	}
	for i := 0; i < 8; i++ {
		x, err := br.read()
		if err != nil {
			panic(err)
		}
		if x {
			t.Errorf("expected all bits to be 0")
		}
	}
	for i := 0; i < 8; i++ {
		x, err := br.read()
		if err != nil {
			panic(err)
		}
		if x != (i%2 == 0) {
			t.Errorf("expected every other bit to be 1")
		}
	}
}

func TestBitReader_readInt(t *testing.T) {
	buf := []byte{0x55} // 01010101
	br := bitReader{r: bytes.NewReader(buf)}
	x, _ := br.readInt(2)
	if x != 1 {
		t.Errorf("expected 1, got %d", x)
	}
	x, _ = br.readInt(4)
	if x != 2 {
		t.Errorf("expected 2, got %d", x)
	}
	x, _ = br.readInt(15)
	if x != 10 {
		t.Errorf("expected 10, got %d", x)
	}

	err := quick.Check(func(x uint) bool {
		buf := make([]byte, 8)
		rand.Read(buf)
		br := bitReader{r: bytes.NewReader(buf)}
		n, err := br.readInt(x)
		if err != nil {
			panic(err)
		}
		return n < x
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
}
