package main

type MoveSorter struct {
	entries [Width]Entry
	size    uint
}

type Entry struct {
	move  uint64
	score int
}

func (ms *MoveSorter) add(move uint64, score int) {
	pos := ms.size
	ms.size++
	// Find pos from max such that the next item has a higher score
	for ; pos > 0 && (ms.entries[pos-1].score > score); pos-- {
		ms.entries[pos] = ms.entries[pos-1]
	}
	ms.entries[pos] = Entry{move, score}
}

func (ms *MoveSorter) getNext() uint64 {
	if ms.size > 0 {
		ms.size--
		return ms.entries[ms.size].move
	} else {
		return 0
	}
}

func (ms MoveSorter) reset() {
	ms.size = 0
}
