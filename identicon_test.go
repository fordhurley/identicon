package identicon

import (
	"math/rand"
	"testing"
	"testing/quick"
)

func TestBitSource_NextBool(t *testing.T) {
	s := BitSource{bytes: []byte{0xFF, 0, 0x55}}
	for i := 0; i < 8; i++ {
		x := s.NextBool()
		if !x {
			t.Errorf("expected all bits to be 1")
		}
	}
	for i := 0; i < 8; i++ {
		x := s.NextBool()
		if x {
			t.Errorf("expected all bits to be 0")
		}
	}
	for i := 0; i < 8; i++ {
		x := s.NextBool()
		if x != (i%2 == 1) {
			t.Errorf("expected every other bit to be 1")
		}
	}
	// Cycle back around:
	for i := 0; i < 8; i++ {
		x := s.NextBool()
		if !x {
			t.Errorf("expected all bits to be 1")
		}
	}
}

func TestBitSource_NextUint(t *testing.T) {
	s := BitSource{bytes: []byte{0x55}} // 01010101

	x := s.NextUint(1)
	if x != 0 {
		t.Errorf("expected 0, got %d", x)
	}
	if s.bitIndex != 0 {
		t.Errorf("expected bitIndex to not be advanced, it is %d", s.bitIndex)
	}

	x = s.NextUint(2)
	if x != 0 {
		t.Errorf("expected 0, got %d", x)
	}
	if s.bitIndex != 1 {
		t.Errorf("expected bitIndex to be advanced to 1, it is %d", s.bitIndex)
	}

	x = s.NextUint(4)
	if x != 2 {
		t.Errorf("expected 2, got %d", x)
	}
	if s.bitIndex != 3 {
		t.Errorf("expected bitIndex to be advanced to 3, it is %d", s.bitIndex)
	}

	x = s.NextUint(15)
	if x != 10 {
		t.Errorf("expected 10, got %d", x)
	}
	if s.bitIndex != 7 {
		t.Errorf("expected bitIndex to be advanced to 7, it is %d", s.bitIndex)
	}

	x = s.NextUint(2)
	if x != 1 {
		t.Errorf("expected 1, got %d", x)
	}
	if s.bitIndex != 0 {
		t.Errorf("expected bitIndex to be reset back to 0, it is %d", s.bitIndex)
	}

	err := quick.Check(func(x uint) bool {
		buf := make([]byte, 8)
		rand.Read(buf)
		s := BitSource{bytes: buf}
		return s.NextUint(x) < x
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
}
