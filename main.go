package main

import "fmt"

type Player uint8

const (
	Player1 Player = 1
	Player2 Player = 2
)

type Column [6]Player

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
func (board *Board) placeStone(column int) {
	if column < 0 || column > 6 {
		panic("Invalid column")
	}
	if board.columns[column].isFull() {
		// print the board, current turn and column
		board.render()
		s := fmt.Sprintf("Column %d is full", column)
		panic(s)
	}
	if board.turn%2 == 0 {
		board.columns[column].placeStone(Player1)
		board.turn++
	} else {
		board.columns[column].placeStone(Player2)
		board.turn++
	}
}

// returns all available columns
func (board *Board) possibleColumns() []int {
	var columns []int
	for i := 0; i < len(board.columns); i++ {
		if !board.columns[i].isFull() {
			columns = append(columns, i)
		}
	}
	return columns
}

func (board *Board) isFull() bool {
	for _, column := range board.columns {
		if !column.isFull() {
			return false
		}
	}
	return true
}

// returns the winner of the game, 0 if no winner
// a winner is defined as a player that has 4 stones in a row, column or any diagonal
func (board *Board) winner() Player {
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

func (board *Board) render() {
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
func (board *Board) takeTurn() {
	if board.isFull() {
		board.render()
		panic("Board is full")
	}
	board.render()
	fmt.Println("Available columns:", board.possibleColumns())
	fmt.Println()
	fmt.Print("Enter column: ")
	var column int
	// check that the user entered a valid column, ask again if not
	for {
		fmt.Scanf("%d", &column)
		if column < 0 || column > 6 {
			fmt.Println("Invalid column")
			fmt.Print("Enter column: ")
		} else if board.columns[column].isFull() {
			fmt.Println("Column is full")
			fmt.Print("Enter column: ")
		} else {
			break
		}
	}
	board.placeStone(column)
}

func main() {
	var board = Board{[7]Column{}, 0}
	for !board.isFull() {
		board.takeTurn()
		// check if there is a winner
		var winner = board.winner()
		if winner != 0 {
			board.render()
			fmt.Println("Winner:", winner)
			break
		}
	}
	board.render()
}

// TODO
