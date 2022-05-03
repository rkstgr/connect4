package main

import (
	"math"
	"sync/atomic"
)

type GamePosition interface {
	isTerminal() bool
	utility() int
	childPositions() []GamePosition
	render()
}


type counter struct {
	count int32
}

// Increment the counter
func (c *counter) inc() {
	atomic.AddInt32(&c.count, 1)
}

func (c *counter) get() int {
	return int(atomic.LoadInt32(&c.count))
}

func minimax(position GamePosition, depth int, maximizingPlayer bool, positionCounter *counter, finished chan bool) int {
	positionCounter.inc()

	if depth == 0 || position.isTerminal() {
		// finished <- true
		return position.utility()
	}
	//finished = make(chan bool)
	if maximizingPlayer {
		maxEval := math.MinInt32 // -infinity
		for _, child := range position.childPositions() {
			eval := minimax(child, depth-1, false, positionCounter, nil)
			// maxEval = max(maxEval, eval)
			if eval > maxEval {
				maxEval = eval
			}
		}
		return maxEval

	} else {
		minEval := math.MaxInt32 // +infinity
		for _, child := range position.childPositions() {
			eval := minimax(child, depth-1, true, positionCounter, nil)
			// minEval = min(minEval, eval)
			if eval < minEval {
				minEval = eval
			}
		}
		return minEval
	}
}