package main

import (
	"testing"
)

func TestWithReadmeExample(t *testing.T) {
	var expected uint64 = 4
	found := SquareSumDifference(2)

	if found != expected {
		t.Errorf("Expected %d but found %d for n=10", expected, found)
	}
}
