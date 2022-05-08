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

// 1 bit at the first cell of every column
var bottomMask = bottom(Width, Height)

// 1 bit for every board cell
var boardMask = bottomMask * ((uint64(1) << Height) - 1)

/*
Board is defined by the bitmask position and the bitmask mask
	* position has a one bit on every cell where the current player has a placed stone
	* mask has a one bit on every cell where some player has placed a stone
	* movesPlayed keeps track of the total number of placed stones
*/
type Board struct {
	position    uint64
	mask        uint64
	movesPlayed int
}

func bottom(width, height int) uint64 {
	if width == 0 {
		return 0
	} else {
		return bottom(width-1, height) | uint64(1<<((width-1)*(height+1)))
	}
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

func (board *Board) play(move uint64) {
	board.position ^= board.mask // invert the stones
	board.mask |= move
	board.movesPlayed++
}

func (board *Board) playColumn(col int) {
	board.play((board.mask + board.bottomMask(col)) & board.columnMask(col))
}

// Plays a sequence of given moves, mainly for setting up a position
// The given numbers stand for the columns starting with 1 with the first column from the left
// Returns the number of moves played, stops the processing when an invalid move is observed:
// - Out of bounds
// - The column is already full
// - Will result to a won position
func (board *Board) playMoves(movesString string) int {
	for i, move := range movesString {
		col := int(move-'0') - 1
		if col < 0 || col > Width || !board.canPlay(col) || board.isWinningMove(col) {
			return i
		}
		board.playColumn(col)
	}
	return len(movesString)
}

func (board *Board) isWinningMove(col int) bool {
	return (board.possible() & board.winningPosition() & board.columnMask(col)) != 0
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

// Bitmask with one at the cell where the next stone can be placed
func (board *Board) possible() uint64 {
	return (board.mask + bottomMask) & boardMask
}

func (board *Board) canWinNext() bool {
	return (board.possible() & board.winningPosition()) != 0
}

func (board *Board) moveScore(move uint64) int {
	return popCount(computeWinningPosition(board.position|move, board.mask))
}

func popCount(a uint64) int {
	c := 0
	for c = 0; a > 0; c++ {
		a &= a - 1
	}
	return c
}

func (board *Board) possibleNonLosingMoves() uint64 {
	possibleMask := board.possible()
	opponentMask := board.opponentWinningPosition()
	forcedMoves := possibleMask & opponentMask
	if forcedMoves > 0 { // We have to prevent the opponent from playing these moves
		if forcedMoves&(forcedMoves-1) > 0 {
			// The opponent has more than one winning move, we can do nothing against it
			return 0
		} else {
			possibleMask = forcedMoves
		}
	}
	return possibleMask & (^(opponentMask >> 1))
}

// Bitmask with one where a placed stone will complete an alignment
func (board *Board) winningPosition() uint64 {
	return computeWinningPosition(board.position, board.mask)
}

func (board *Board) opponentWinningPosition() uint64 {
	return computeWinningPosition(board.position^board.mask, board.mask)
}

// Computes all the bits in the position that complete an alignment
func computeWinningPosition(position, mask uint64) uint64 {
	// vertical
	var r = (position << 1) & (position << 2) & (position << 3) // three consecutive stones below

	// horizontal
	var p uint64 = (position << (Height + 1)) & (position << (2 * (Height + 1))) // two stones on the left
	r |= p & (position << (3 * (Height + 1)))                                    // the third on the left
	r |= p & (position >> (Height + 1))                                          // one on the right
	p = (position >> (Height + 1)) & (position >> (2 * (Height + 1)))            // two stones on the right
	r |= p & (position >> (3 * (Height + 1)))                                    // the third on the right
	r |= p & (position << (Height + 1))                                          // one on the left

	// diagonal 1 (bottom left -> top right)
	p = (position << (Height + 2)) & (position << (2 * (Height + 2))) // two stones bottom left
	r |= p & (position << (3 * (Height + 2)))                         // third stone bottom left
	r |= p & (position >> (Height + 2))                               // one stone top right
	p = (position >> (Height + 2)) & (position >> (2 * (Height + 2))) // two stones top right
	r |= p & (position >> (3 * (Height + 2)))                         // third stone top right
	r |= p & (position << (Height + 2))                               // one stone bottom left

	// diagonal 2 (top left -> bottom right)
	p = (position << Height) & (position << (2 * Height)) // two stones top left
	r |= p & (position << (3 * Height))                   // third stone top left
	r |= p & (position >> Height)                         // one stone bottom right
	p = (position >> Height) & (position >> (2 * Height)) // two stones bottom right
	r |= p & (position >> (3 * Height))                   // third stone bottom right
	r |= p & (position << Height)                         // one stone top left

	return r & (boardMask ^ mask)
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

func newBoard(movesString string) Board {
	board := Board{0, 0, 0}
	board.playMoves(movesString)
	return board
}

func main() {
	board := newBoard("271713432331713132")
	counter := Counter{}
	eval := solve(board, &counter)
	print(eval)
}
