package main

import (
	"fmt"
)

type Player uint8

//goland:noinspection GoUnusedConst
const (
	Player1  Player = 1
	Player2  Player = 2
	Height   int    = 6
	Width    int    = 7
	MaxScore        = Height*Width/2 - 3
	MinScore        = -Height*Width/2 + 3
)

type Board struct {
	position    uint64
	mask        uint64
	bottom      uint64
	movesPlayed int
}

func (board *Board) minScore() int {
	return -(Width*Height - board.movesPlayed) / 2
}

func (board *Board) maxScore() int {
	return (Width*Height + 1 - board.movesPlayed) / 2
}

// Return true if there is no stone in the top cell of the column
func (board *Board) canPlay(col int) bool {
	return (board.mask & board.topMask(col)) == 0
}

func (board *Board) play(col int) {
	board.position ^= board.mask // invert the stones
	board.mask |= board.mask + board.bottomMask(col)
	board.movesPlayed++
}

// Plays a sequence of given moves, mainly for setting up a position
// The given numbers stand for the columns starting with 1 with the first column from the left
// Returns the number of moves played, stops the processing when an invalid move is observed:
// - Out of bounds
// - The column is already full
// - Will result to a won position
func (board *Board) playMoves(movesString string) int {
	for i, move := range movesString {
		col := int(move - '1')
		if col < 0 || col > Width || !board.canPlay(col) || board.isWinningMove(col) {
			return i
		}
		board.play(col)
	}
	return len(movesString)
}

func (board *Board) isWinningMove(col int) bool {
	pos := board.position
	pos |= (board.mask + board.bottomMask(col)) & board.columnMask(col)
	return alignment(pos)
}

func (board *Board) key() uint64 {
	return board.position + board.mask
}

// Return a bitmask with a single one in the top cell of the column
func (board *Board) topMask(col int) uint64 {
	return uint64((1 << (Height - 1)) << (col * (Height + 1)))
}

// Return a bitmask with a single one in the first cell of the column
func (board *Board) bottomMask(col int) uint64 {
	return uint64(1 << (col * (Height + 1)))
}

// Return a bitmask with 1 for every cell of the column
func (board *Board) columnMask(col int) uint64 {
	return uint64(((1 << Height) - 1) << (col * (Height + 1)))
}

func (board *Board) render() {
	// Loop over every bit from top left to right, downwards
	for j := Height - 1; j >= 0; j-- {
		for i := 0; i < Width; i++ {
			bit := uint64(1 << (i*(Height+1) + j))
			// Check if empty
			if (board.mask & bit) == 0 {
				fmt.Print("🔘") // Empty
			} else {
				if (board.position & bit) != 0 {
					fmt.Print("🔴") // Player
				} else {
					fmt.Print("🟡") // Opponent
				}
			}
		}
		// End of line
		fmt.Println("")
	}
	// End
	fmt.Println("")
}

// Return true if there is an alignment of four stones in any direction for the current player position
func alignment(pos uint64) bool {
	// vertical
	var m uint64 = pos & (pos >> 1)
	if (m & (m >> 2)) != 0 {
		return true
	}

	// horizontal
	m = pos & (pos >> (Height + 1))
	if (m & (m >> (2 * (Height + 1)))) != 0 {
		return true
	}

	// diagonal left up
	m = pos & (pos >> Height)
	if (m & (m >> (2 * Height))) != 0 {
		return true
	}

	// diagonal left down
	m = pos & (pos >> (Height + 2))
	if (m & (m >> (2 * (Height + 2)))) != 0 {
		return true
	}

	return false
}

func newBoard(movesString string) Board {
	board := Board{0, 0, 0, 0}
	board.playMoves(movesString)
	return board
}

func main() {
	board := newBoard("7422341735647741166133573473242566")
	board.render()
}
