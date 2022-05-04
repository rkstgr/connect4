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

func negamax(board Board, alpha, beta int, counter *Counter) int {
	counter.inc()
	var i = 0

	// Check if the game is drawn
	if board.movesPlayed == Height*Width {
		return 0
	}

	// Check if the current player has a winning move
	for ; i < Width; i++ {
		if board.canPlay(i) && board.isWinningMove(i) {
			return (Width*Height + 1 - board.movesPlayed) / 2
		}
	}

	// Since we cannot win with the next move, the best score is if we move with our second next move
	// Which is the score from above + 1
	max := (Width*Height - 1 - board.movesPlayed) / 2 // upper bound of our score as we cannot win immediately
	if beta > max {
		beta = max // there is no need to keep beta above our max possible score.

		if alpha >= beta {
			// There is another path with a higher score, no need to continue
			return beta // prune the exploration if the [alpha;beta] window is empty.
		}
	}

	var moves = [7]int{3, 2, 4, 1, 5, 0, 6}
	i = 0
	for ; i < 7; i++ {
		move := moves[i]
		if board.canPlay(move) {
			nextBoard := board
			nextBoard.playMove(move)
			score := -negamax(nextBoard, -beta, -alpha, counter)
			if score >= beta {
				// There is better path for the opponent; Prune
				return score
			}
			// Update alpha if possible
			if score > alpha {
				alpha = score
			}
		}
	}
	return alpha
}
