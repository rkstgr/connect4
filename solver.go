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
	if alpha >= beta {
		panic("Negamax alpha should be smaller beta")
	}
	if board.canWinNext() {
		panic("Negamax should not be called on a board with a winning move")
	}
	counter.inc()

	possibleNonLosingMoves := board.possibleNonLosingMoves()
	if possibleNonLosingMoves == 0 {
		// We have no non-losing moves; opponent wins next move
		return -(Height*Width - board.movesPlayed) / 2
	}

	// Check for drawn game (we look two steps ahead)
	if board.movesPlayed >= ((Height * Width) - 2) {
		return 0
	}

	min := -(Height*Width - 2 - board.movesPlayed) / 2
	if alpha < min {
		alpha = min
		if alpha >= beta {
			return alpha
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

	moves := MoveSorter{}
	var moveOrder = []int{3, 2, 4, 1, 5, 0, 6}

	for i := 0; i < Width; i++ {
		if move := possibleNonLosingMoves & board.columnMask(moveOrder[i]); move > 0 {
			moves.add(move, board.moveScore(move))
		}
	}

	for next := moves.getNext(); next > 0; next = moves.getNext() {
		nextBoard := board
		nextBoard.play(next)
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
	transpositionTable.set(board.key(), int8(alpha-MinScore+1))
	return alpha
}

func solve(board Board, counter *Counter) int {
	if board.canWinNext() {
		return (Height*Width + 1 - board.movesPlayed) / 2
	}
	min := -(Height*Width - board.movesPlayed) / 2
	max := (Height*Width + 1 - board.movesPlayed) / 2

	for min < max {
		med := min + (max-min)/2
		if med <= 0 && min/2 < med {
			med = min / 2
		} else if med >= 0 && max/2 > med {
			med = max / 2
		}
		r := negamax(board, med, med+1, counter) // use a null depth window to know if the actual score is greater or smaller than med
		if r <= med {
			max = r
		} else {
			min = r
		}
	}
	return min
}
