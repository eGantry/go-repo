package main

import (
	"fmt"
)

const size = 9

var board [size][size]string
var currentPlayer = "B" // B = Black, W = White

// Initialize the board
func initBoard() {
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			board[i][j] = "."
		}
	}
}

// Print the board
func printBoard() {
	fmt.Println("   0 1 2 3 4 5 6 7 8")
	for i := 0; i < size; i++ {
		fmt.Printf("%d  ", i)
		for j := 0; j < size; j++ {
			fmt.Print(board[i][j], " ")
		}
		fmt.Println()
	}
}

// Place a stone on the board
func placeStone(x, y int) bool {
	if x < 0 || x >= size || y < 0 || y >= size {
		fmt.Println("Invalid move! Out of bounds.")
		return false
	}
	if board[x][y] != "." {
		fmt.Println("Invalid move! Spot is already taken.")
		return false
	}
	board[x][y] = currentPlayer
	return true
}

// Switch player
func switchPlayer() {
	if currentPlayer == "B" {
		currentPlayer = "W"
	} else {
		currentPlayer = "B"
	}
}

func main() {
	initBoard()
	fmt.Println("Simple Go Game (9x9) - No captures yet")
	printBoard()

	for {
		var x, y int
		fmt.Printf("Player %s, enter row and column (e.g., '4 5'): ", currentPlayer)
		_, err := fmt.Scan(&x, &y)
		if err != nil {
			fmt.Println("Invalid input! Enter two numbers separated by a space.")
			continue
		}

		if placeStone(x, y) {
			printBoard()
			switchPlayer()
		}
	}
}
