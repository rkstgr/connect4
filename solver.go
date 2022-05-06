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

// Counter counting visited positions
type Counter struct {
	mu    sync.Mutex
	count int32
}

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

// TranspositionTable Caches values of game positions
type TranspositionTable struct {
	table map[uint64]int8
}

func NewTranspositionTable() *TranspositionTable {
	return &TranspositionTable{
		table: make(map[uint64]int8),
	}
}

func (t *TranspositionTable) get(key uint64) int8 {
	// check if columnKey exists is not return 0
	if value, ok := t.table[key]; ok {
		return value
	} else {
		return 0
	}
}

func (t *TranspositionTable) set(key uint64, value int8) {
	t.table[key] = value
}

var transpositionTable = NewTranspositionTable()

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
	if val := transpositionTable.get(board.key()); val != 0 {
		max = int(val) + MinScore - 1
	}
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
			nextBoard.play(move)
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
	transpositionTable.set(board.key(), int8(alpha-MinScore+1))
	return alpha
}

func solve(board Board, counter *Counter) int {
	min := board.minScore()
	max := board.maxScore()
	for min < max {
		middle := min + (max-min)/2
		if middle <= 0 && min/2 < middle {
			middle = min / 2
		} else if middle >= 0 && middle < max/2 {
			middle = max / 2
		}
		r := negamax(board, middle, middle+1, counter)
		if r <= middle {
			max = r
		} else {
			min = r
		}
	}
	return min
}
