package main

import (
	"fmt"
)

type Player uint8

const (
	Player1 Player = 1
	Player2 Player = 2
	Height  int    = 6
	Width   int    = 7
)

type Column [Height]Player

type Board struct {
	columns     [Width]Column
	heights     [Width]int
	position    string
	movesPlayed int
}

// COLUMN

// Place a stone at the end of the given column
func (c *Column) placeStone(player Player) {
	for i := 0; i < len(c); i++ {
		if c[i] == 0 {
			c[i] = player
			return
		}
	}
}

// Check if the given column is full
// check if the last element is not 0
func (c *Column) isFull() bool {
	return c[len(c)-1] != 0
}

// BOARD

// check if the move is possible
func (board Board) canPlay(move int) bool {
	return move >= 0 && move < len(board.columns) && board.heights[move] < Height
}

func (board *Board) makeMove(move int, player Player) {
	board.columns[move].placeStone(player)
	board.heights[move]++
	board.movesPlayed++
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
	if board.columns[move].isFull() {
		// print the board, current turn and move
		board.render()
		s := fmt.Sprintf("Column %d is full", move)
		panic(s)
	}
	if board.movesPlayed%2 == 0 {
		board.columns[move].placeStone(Player1)
	} else {
		board.columns[move].placeStone(Player2)
		// board.position += fmt.Sprintf("%d", move+1)
	}
	board.heights[move]++
	board.movesPlayed++
}

// Returns all available columns
// If a heuristic is given it will return the moves sorted by their heuristic value
func (board *Board) possibleMoves() []int {
	var moves [7]int
	movesFound := 0
	moveOrder := []int{3, 2, 4, 1, 5, 0, 6} // move heuristic: center best, corners worst
	for _, move := range moveOrder {
		if !board.columns[move].isFull() {
			moves[movesFound] = move
			movesFound++
		}
	}
	return moves[:movesFound]
}

func (board *Board) stonesPlacedBy(player Player) int {
	// sum the number of stones in each column
	var sum int
	for _, column := range board.columns {
		for _, stone := range column {
			if stone == player {
				sum++
			}
		}
	}
	return sum
}

// Returns the player who has to make the next move
func (board *Board) currentPlayer() Player {
	if board.movesPlayed%2 == 0 {
		return Player1
	}
	return Player2
}

// Returns true if there is no move possible
// Because all columns are full
func (board *Board) isFull() bool {
	for _, column := range board.columns {
		if !column.isFull() {
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
func makeManualMoves(board Board) Board {
	if board.isFull() {
		board.render()
		panic("Board is full")
	}
	var move int
	// check that the user entered a valid move, ask again if not
	for !board.isTerminal() {
		board.render()
		fmt.Println("Available columns:", board.possibleMoves())
		fmt.Println()
		fmt.Print("Enter move: ")
		for {
			fmt.Scanf("%d", &move)
			if move == -1 {
				return board
			}
			if move < 0 || int(move) >= len(board.columns) {
				fmt.Println("Invalid move")
				fmt.Print("Enter move: ")
			} else if board.columns[move].isFull() {
				fmt.Println("Column is full")
				fmt.Print("Enter move: ")
			} else {
				break
			}
		}
		// board = playMove(board, move)
	}
	return board
}

// GAME

func (board Board) isTerminal() bool {
	return board.isFull() || board.wonBy() != 0
}

// Score of the current position relative to the current player
// A positive score means the current player has a winning position (won the game)
// A negative score means the current player has a losing position (lost the game)
// A score of 0 means the every move will lead to a draw
func (board Board) positionScore() int {
	if board.isTerminal() {
		winningPlayer := board.wonBy()

		// Player 1 has a winning position
		if winningPlayer == Player1 {
			return 22 - board.stonesPlacedBy(Player1)
		} else if winningPlayer == Player2 {
			// Current player has a losing position
			return -22 + board.stonesPlacedBy(Player2)
		} else {
			return 0
		}
	} else {
		// The current position is not terminal
		return 0
	}
}

func (board Board) getPosition() string {
	return ""
}

func (board Board) evaluate() int {
	positionsVisited := Counter{}
	return negamax(board, &positionsVisited)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// returns the heuristic value of the board
// the value is highest in the center and decreases linearly to the edges
func centerHeuristic(_ *Board, move int) int {
	return 3 - abs(move-3)
}

func (board Board) childPositions() []Board {
	var children []Board
	for _, move := range board.possibleMoves() {
		child := board
		child.playMove(move)
		children = append(children, child)
	}
	return children
}

// movesString is a string of integers representing the moves read from left to right
// e.g "1234" means the first move is column 1, the second move is column 2, etc.
func createBoard(position string) Board {
	var board Board
	for i := 0; i < len(position); i++ {
		move := int((position[i] - '0') - 1)
		board.playMove(move)
	}
	return board
}

func main() {
	//var board = Board{[7]Column{}, 0}
	board := createBoard("67635256351344534443614126713657127")
	fmt.Println(board.position)
	fmt.Println("Current Player:", board.currentPlayer())
	fmt.Println(board.positionScore())
	board.render()
	counter := Counter{}
	eval := negamax(board, &counter)

	fmt.Println("Evaluation:", eval)
	fmt.Println("Positions visited:", counter.count)
}
