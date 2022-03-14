package main

import (
	"math"
)

// https://codereview.stackexchange.com/questions/868/calculating-entropy-of-a-string
func calcShanonEntropy(s string) float64 {
	c2n := make(map[rune]int)

	rs := []rune(s)
	for _, r := range rs {
		cnt, exists := c2n[r]
		if exists {
			c2n[r] = cnt + 1
		} else {
			c2n[r] = 1
		}
	}

	result := 0.0
	slen := len(rs)

	for _, v := range c2n {
		freq := float64(v) / float64(slen)

		result -= freq * (math.Log(freq) / math.Log(2.0))
	}

	return result
}
