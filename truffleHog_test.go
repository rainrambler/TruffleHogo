package main

import (
	"testing"
)

func TestParseTextLine1(t *testing.T) {
	s := "1abisdfse3355ve4e34 dfslfusdnvdfnjd6 zsdjfnej aaaaaaaaaaaaaaaaa"

	res := parseTextLine(1, s)
	expected := 2

	if res != expected {
		t.Errorf("Result: %v, want: %v", res, expected)
	}
}
