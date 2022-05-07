package main

import (
	"testing"
)

func TestMoveSorter(t *testing.T) {
	moves := MoveSorter{}

	moves.add(1, 1)
	moves.add(3, 3)
	moves.add(2, 2)
	moves.add(4, 2)

	sorted := []uint64{}

	for m := moves.getNext(); m > 0; m = moves.getNext() {
		print(m, " ")
		sorted = append(sorted, m)
	}
	y := []uint64{3, 2, 4, 1}
	for i, m := range sorted {
		if m != y[i] {
			t.Fatalf("Fail")
		}
	}

}
