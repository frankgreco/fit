package main

import (
	"fmt"
	"testing"
)

func TestExplodeByte(t *testing.T) {
	for _, test := range []struct {
		in  uint8
		out string
	}{
		{0, "00000000"},
		{255, "11111111"},
		{1, "00000001"},
		{254, "11111110"},
		{32, "00100000"},
		{132, "10000100"},
	} {
		if ExplodeByte(test.in) != test.out {
			fmt.Println(ExplodeByte(test.in))
			t.Fail()
		}
	}
}
