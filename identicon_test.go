package identicon

import (
	"math/rand"
	"testing"
	"testing/quick"
)

func TestBitSource_nextBool(t *testing.T) {
	s := bitSource{bytes: []byte{0xFF, 0, 0x55}}
	for i := 0; i < 8; i++ {
		x := s.nextBool()
		if !x {
			t.Errorf("expected all bits to be 1")
		}
	}
	for i := 0; i < 8; i++ {
		x := s.nextBool()
		if x {
			t.Errorf("expected all bits to be 0")
		}
	}
	for i := 0; i < 8; i++ {
		x := s.nextBool()
		if x != (i%2 == 0) {
			t.Errorf("expected every other bit to be 1")
		}
	}
	// Cycle back around:
	for i := 0; i < 8; i++ {
		x := s.nextBool()
		if !x {
			t.Errorf("expected all bits to be 1")
		}
	}
}

func TestBitSource_nextUint(t *testing.T) {
	s := bitSource{bytes: []byte{0x55}} // 01010101
	x := s.nextUint(2)
	if x != 1 {
		t.Errorf("expected 1, got %d", x)
	}
	x = s.nextUint(4)
	if x != 2 {
		t.Errorf("expected 2, got %d", x)
	}
	x = s.nextUint(15)
	if x != 10 {
		t.Errorf("expected 10, got %d", x)
	}

	err := quick.Check(func(x uint) bool {
		buf := make([]byte, 8)
		rand.Read(buf)
		s := bitSource{bytes: buf}
		return s.nextUint(x) < x
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
}
