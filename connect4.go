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
	MaxScore int    = Height*Width/2 - 3
	MinScore int    = -Height*Width/2 + 3
)

type Column [Height]Player

type Board struct {
	columns     [Width]Column
	heights     [Width]int
	position    string
	movesPlayed int
}

// BOARD

func (board Board) canPlay(move int) bool {
	return move >= 0 && move < len(board.columns) && board.heights[move] < Height
}

func (board *Board) undoMove(move int) {
	for i := len(board.columns[move]) - 1; i >= 0; i-- {
		if board.columns[move][i] != 0 {
			board.columns[move][i] = 0
			board.heights[move]--
			return
		}
	}
}

func (board Board) isWinningMove(move int) bool {

	currentPlayer := board.currentPlayer()

	// check for vertical alignments
	if board.heights[move] >= 3 &&
		board.columns[move][board.heights[move]-1] == currentPlayer &&
		board.columns[move][board.heights[move]-2] == currentPlayer &&
		board.columns[move][board.heights[move]-3] == currentPlayer {
		return true
	}

	for dy := -1; dy <= 1; dy++ { // Iterate on horizontal (dy = 0) or two diagonal directions (dy = -1 or dy = 1).
		nb := 0                          // counter of the number of stones of current player surronding the played stone in tested direction.
		for dx := -1; dx <= 1; dx += 2 { // count continuous stones of current player on the left, then right of the played column.
			x := move + dx
			y := board.heights[move] + dx*dy
			for ; nb < 3 && x >= 0 && x < Width && y >= 0 && y < Height && board.columns[x][y] == currentPlayer; nb++ {
				x += dx
				y += dy * dx
			}
		}
		if nb == 3 {
			return true
		}
		// there is an aligment if at least 3 other stones of the current user
		// are surronding the played stone in the tested direction.
	}
	return false
}

// places a stone into the given column, panic if the given column is already full
// the color of the next stone is based on the current turn given that player 1 always starts
func (board *Board) playMove(move int) {
	if move < 0 || move >= len(board.columns) {
		panic("Invalid move")
	}
	// check if column is full
	currentHeight := board.heights[move]
	if currentHeight >= Height {
		s := fmt.Sprintf("Column %d is full", move)
		panic(s)
	}
	if board.movesPlayed%2 == 0 {
		board.columns[move][currentHeight] = Player1
	} else {
		board.columns[move][currentHeight] = Player2
	}
	board.heights[move]++
	board.movesPlayed++
}

func (board *Board) currentPlayer() Player {
	if board.movesPlayed%2 == 0 {
		return Player1
	}
	return Player2
}

func (board *Board) isFull() bool {
	for _, height := range board.heights {
		if height < Height {
			return false
		}
	}
	return true
}

// Checks if the current position is already won and return the winning player
// If not return 0
// A winner is defined as a player that has 4 stones in a row, column or any diagonal
func (board *Board) wonBy() Player {
	for i := 0; i < len(board.columns); i++ {
		for j := 0; j < len(board.columns[i]); j++ {
			if winningPlayer := board.columns[i][j]; winningPlayer != 0 {
				// check row
				if i+3 < len(board.columns) && board.columns[i+1][j] == winningPlayer && board.columns[i+2][j] == winningPlayer && board.columns[i+3][j] == winningPlayer {
					return winningPlayer
				}
				// check column
				if j+3 < len(board.columns[i]) && board.columns[i][j+1] == winningPlayer && board.columns[i][j+2] == winningPlayer && board.columns[i][j+3] == winningPlayer {
					return winningPlayer
				}
				// check diagonal right up
				if i+3 < len(board.columns) && j+3 < len(board.columns[i]) && board.columns[i+1][j+1] == winningPlayer && board.columns[i+2][j+2] == winningPlayer && board.columns[i+3][j+3] == winningPlayer {
					return winningPlayer
				}
				// check diagonal right down
				if i+3 < len(board.columns) && j-3 >= 0 && board.columns[i+1][j-1] == winningPlayer && board.columns[i+2][j-2] == winningPlayer && board.columns[i+3][j-3] == winningPlayer {
					return winningPlayer
				}
			}
		}
	}
	return 0
}

// Prints the board to the console
func (board Board) render() {
	for j := len(board.columns[0]) - 1; j >= 0; j-- {
		for _, column := range board.columns {
			position := column[j]
			if position == 0 {
				print("🔘")
			} else if position == Player1 {
				print("🟢")
			} else if position == Player2 {
				print("🔴")
			}
			print(" ")
		}
		println()
	}
}

// Take turn manually, panics if board is full
// Render the board and ask the user for a column
//goland:noinspection GoUnusedFunction
func (board Board) makeManualMove() {
	// Check if a move is possible
	if board.isFull() {
		fmt.Println("The board is full")
		return
	}

	board.render()
	fmt.Println("Please enter a valid move: ")
	var move int
	_, err := fmt.Scanf("%d", &move)
	if err != nil {
		fmt.Println("Error parsing the entered move")
		return
	}
	if board.canPlay(move) {
		board.playMove(move)
	} else {
		fmt.Println("Invalid move")
	}
}

func (board Board) negamaxScore() int {
	positionsVisited := Counter{}
	return negamax(board, -1000, 1000, &positionsVisited)
}

// movesString is a string of integers representing the moves read from left to right
// e.g "1234" means the first move is column 1, the second move is column 2, etc.
func createBoard(positionString string) Board {
	var board Board
	for i := 0; i < len(positionString); i++ {
		move := int((positionString[i] - '0') - 1)
		board.playMove(move)
	}
	return board
}

func (board Board) columnKey(column int) uint8 {
	var key uint8
	columnHeight := board.heights[column]
	if columnHeight == 0 {
		return 0
	}
	key = 1<<(columnHeight) - 1
	for i := columnHeight - 1; i >= 0; i-- {
		if board.columns[column][i] == Player2 {
			key += 1 << (columnHeight - i - 1)
		}
	}
	return key
}

// Get the key for the current position
// Symmetric positions are treated as equal
func (board *Board) key() uint64 {
	columnKeys := [Width]uint8{}
	for i := 0; i < Width; i++ {
		columnKeys[i] = board.columnKey(i)
	}
	var leftFirst = false
	if columnKeys[0]+columnKeys[1]+columnKeys[2] > columnKeys[4]+columnKeys[5]+columnKeys[6] {
		leftFirst = true
	}
	var key uint64
	if leftFirst {
		for i := 0; i < Width; i++ {
			key = key << 8
			key += uint64(columnKeys[i])
		}
	} else {
		for i := Width - 1; i >= 0; i-- {
			key = key << 8
			key += uint64(columnKeys[i])
		}
	}
	return key
}

func main() {
	//var board = Board{[7]Column{}, 0}
	board := createBoard("34225")
	fmt.Println("Current Player:", board.currentPlayer())
	board.render()

	counter := Counter{}
	eval := negamax(board, -1000, 1000, &counter)

	fmt.Println("Evaluation:", eval)
	fmt.Println("Positions visited:", counter.count)
}
