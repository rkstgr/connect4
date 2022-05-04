package main

import (
	"math"
	"sync/atomic"
)

type GamePosition interface {
	isTerminal() bool
	utility() int
	childPositions() []GamePosition
	getPosition() string
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

func negamax(board Board) (int, int) {
	// Drawn game
	positionsCount := 0
	if board.isTerminal() && board.wonBy() == 0 {
		return 0, positionsCount
	}

	for _, move := range board.possibleMoves(centerHeuristic) {
		positionsCount++
		if board.isWinningMove(move) {
			return (Width*Height + 1 - board.stonesPlaced()) / 2, positionsCount
		}
	}

	var bestScore = -Width * Height
	for _, move := range board.possibleMoves(centerHeuristic) {
		nextBoard := playMove(board, move)
		score, otherCount := negamax(nextBoard)
		score = -score
		positionsCount += otherCount
		if score > bestScore {
			bestScore = score
		}
	}
	return bestScore, positionsCount
}

func minimax(board Board, depth int, maximizingPlayer bool, positionCounter *counter) int {
	positionCounter.inc()

	if depth == 0 || board.isTerminal() {
		return board.positionScore()
	}
	if maximizingPlayer {
		maxEval := math.MinInt32 // -infinity
		for _, nextBoard := range board.childPositions() {
			eval := minimax(nextBoard, depth-1, false, positionCounter)
			// maxEval = max(maxEval, eval)
			if eval > maxEval {
				maxEval = eval
			}
		}
		return maxEval

	} else {
		minEval := math.MaxInt32 // +infinity
		for _, child := range board.childPositions() {
			eval := minimax(child, depth-1, true, positionCounter)
			// minEval = min(minEval, eval)
			if eval < minEval {
				minEval = eval
			}
		}
		return minEval
	}
}
