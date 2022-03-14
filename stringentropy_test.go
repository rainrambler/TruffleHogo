package main

import (
	"math"
	"testing"
)

func TestCalcShanonEntropy1(t *testing.T) {
	s := "aaaa"

	res := calcShanonEntropy(s)
	expected := 0.0

	if res != expected {
		t.Errorf("Result: %v, want: %v", res, expected)
	}
}

func TestCalcShanonEntropy2(t *testing.T) {
	s := "123123123123"

	res := calcShanonEntropy(s)
	expected := 1.5

	if !isInRange(res, expected) {
		t.Errorf("Result: %v, want: %v", res, expected)
	}
}

func TestCalcShanonEntropy3(t *testing.T) {
	s := "abcdefghijklmnop"

	res := calcShanonEntropy(s)
	expected := 4.0

	if !isInRange(res, expected) {
		t.Errorf("Result: %v, want: %v", res, expected)
	}
}

func TestCalcShanonEntropy4(t *testing.T) {
	s := "YWJjZGVm==" // BASE64 of "abcdef"

	res := calcShanonEntropy(s)
	expected := 3.1

	if !isInRange(res, expected) {
		t.Errorf("Result: %v, want: %v", res, expected)
	}
}

func isInRange(v, base float64) bool {
	return math.Abs(v-base) < 0.1
}
