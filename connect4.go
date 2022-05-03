package main

import (
	"fmt"
	"sort"
)

type Player uint8

const (
	Player1 Player = 1
	Player2 Player = 2
)

type Column [6]Player

type Move int

type Board struct {
	columns [7]Column
	turn    int
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
func (c *Column) isFull() bool {
	for _, stone := range c {
		if stone == 0 {
			return false
		}
	}
	return true
}

// BOARD

// places a stone into the given column, panic if the given column is already full
// the color of the next stone is based on the current turn given that player 1 always starts
func playMove(board Board, move Move) Board {
	if move < 0 || int(move) >= len(board.columns) {
		panic("Invalid move")
	}
	if board.columns[move].isFull() {
		// print the board, current turn and move
		board.render()
		s := fmt.Sprintf("Column %d is full", move)
		panic(s)
	}
	if board.turn%2 == 0 {
		board.columns[move].placeStone(Player1)
		board.turn++
	} else {
		board.columns[move].placeStone(Player2)
		board.turn++
	}
	return board
}

// returns all available columns
// if a moveHeuristic is given it will return the moves sorted by their moveHeuristic
func (board *Board) possibleMoves(moveHeuristic func(b *Board, m Move) int) []Move {
	var moves []Move
	for i := 0; i < len(board.columns); i++ {
		if !board.columns[i].isFull() {
			moves = append(moves, Move(i))
		}
	}
	if moveHeuristic != nil {
		sort.Slice(moves, func(i, j int) bool {
			return moveHeuristic(board, moves[i]) > moveHeuristic(board, moves[j])
		})
	}
	return moves
}

func (board *Board) whichTurn() Player {
	if board.turn%2 == 0 {
		return Player1
	}
	return Player2
}

func (board *Board) isFull() bool {
	for _, column := range board.columns {
		if !column.isFull() {
			return false
		}
	}
	return true
}

// returns the winner of the game, 0 if there is no winner or the game is drawn
// a winner is defined as a player that has 4 stones in a row, column or any diagonal
func (board *Board) evaluate() Player {
	for i := 0; i < len(board.columns); i++ {
		for j := 0; j < len(board.columns[i]); j++ {
			if player := board.columns[i][j]; player != 0 {
				// check row
				if i+3 < len(board.columns) && board.columns[i+1][j] == player && board.columns[i+2][j] == player && board.columns[i+3][j] == player {
					return player
				}
				// check column
				if j+3 < len(board.columns[i]) && board.columns[i][j+1] == player && board.columns[i][j+2] == player && board.columns[i][j+3] == player {
					return player
				}
				// check diagonal right up
				if i+3 < len(board.columns) && j+3 < len(board.columns[i]) && board.columns[i+1][j+1] == player && board.columns[i+2][j+2] == player && board.columns[i+3][j+3] == player {
					return player
				}
				// check diagonal right down
				if i+3 < len(board.columns) && j-3 >= 0 && board.columns[i+1][j-1] == player && board.columns[i+2][j-2] == player && board.columns[i+3][j-3] == player {
					return player
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
	var move Move
	// check that the user entered a valid move, ask again if not
	for !board.isTerminal() {
		board.render()
		fmt.Println("Available columns:", board.possibleMoves(nil))
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
		board = playMove(board, move)
	}
	return board
}

// GAME

func (board Board) isTerminal() bool {
	return board.evaluate() != 0 || board.isFull()
}

func (board Board) utility() int {
	if board.evaluate() == Player1 {
		return 1
	} else if board.evaluate() == Player2 {
		return -1
	}
	return 0
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// returns the heuristic value of the board
// the value is highest in the center and decreases linearly to the edges
func centerHeuristic(_ *Board, move Move) int {
	return 3 - abs(int(move-3))
}

func (board Board) childPositions() []GamePosition {
	var children []GamePosition
	for _, move := range board.possibleMoves(centerHeuristic) {
		child := playMove(board, move)
		children = append(children, child)
	}
	return children
}

// movesString is a string of integers representing the moves read from left to right
// e.g "1234" means the first move is column 1, the second move is column 2, etc.
func createBoard(position string) Board {
	var board Board
	for i := 0; i < len(position); i++ {
		move := Move((position[i] - '0') - 1)
		board = playMove(board, move)
	}
	return board
}

func main() {
	//var board = Board{[7]Column{}, 0}
	board := createBoard("2252576253462244111563365343671351441")
	board.render()

	positionsVisited := counter{}
	finished := make(chan bool)
	eval := minimax(board, 20, board.whichTurn() == Player1, &positionsVisited, finished)

	fmt.Println("Positions visited:", positionsVisited.count)
	fmt.Println("Evaluation:", eval)
}
