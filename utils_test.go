package gobencode

import (
	"testing"
)

func TestByteToInteger(t *testing.T) {
	for b := byte(0); b < '0'; b++ {
		if _, ok := byteToInteger(b); ok {
			t.Fatalf("byte [\\%x] is not an integer, output: [%v] | expected: [%v]", b, ok, false)
		}
	}
	i := 0
	for b := byte('0'); b <= '9'; b++ {
		integer, ok := byteToInteger(b)

		if !ok {
			t.Fatalf("byte [\\%x] is an integer: output [%v] | expected [%v]", b, ok, true)
		}
		if integer != i {
			t.Fatalf("bad integer converion: output [%v] | expected [%v]", integer, i)
		}
		i++
	}
	for b := byte('9' + 1); b < ^byte(0); b++ {
		if _, ok := byteToInteger(b); ok {
			t.Fatalf("byte [\\%x] is not an integer: output [%v] | expected [%v]", b, ok, false)
		}
	}

}
