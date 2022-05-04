package main

import (
	"sync"
)

type GamePosition interface {
	isTerminal() bool
	utility() int
	childPositions() []GamePosition
	getPosition() string
	render()
}

type Counter struct {
	mu    sync.Mutex
	count int32
}

// Increment the Counter
func (c *Counter) inc() {
	c.mu.Lock()
	c.count++
	c.mu.Unlock()
}

func (c *Counter) get() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return int(c.count)
}

func negamax(board Board, counter *Counter) int {
	counter.inc()
	var i = 0
	// Drawn game
	if board.movesPlayed == Height*Width {
		return 0
	}

	for ; i < Width; i++ {
		if board.canPlay(i) && board.isWinningMove(i) {
			return (Width*Height + 1 - board.movesPlayed) / 2
		}
	}

	var bestScore = -Width * Height
	var moves = [7]int{3, 2, 4, 1, 5, 0, 6}
	i = 0
	for ; i < 7; i++ {
		move := moves[i]
		if board.canPlay(move) {
			nextBoard := board
			nextBoard.playMove(move)
			score := -negamax(nextBoard, counter)
			if score > bestScore {
				bestScore = score
			}
		}
	}
	return bestScore
}
