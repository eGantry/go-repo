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

// Place a stone and check for captures
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
	checkForCaptures()
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

// Check for captured stones and remove them
func checkForCaptures() {
	opponent := "B"
	if currentPlayer == "B" {
		opponent = "W"
	}

	// Scan board for groups with no liberties
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			if board[i][j] == opponent {
				visited := make(map[[2]int]bool)
				if !hasLiberty(i, j, opponent, visited) {
					removeCapturedGroup(visited)
				}
			}
		}
	}
}

// Check if a group has a liberty (empty space nearby)
func hasLiberty(x, y int, color string, visited map[[2]int]bool) bool {
	if x < 0 || x >= size || y < 0 || y >= size {
		return false
	}
	if board[x][y] == "." {
		return true
	}
	if board[x][y] != color || visited[[2]int{x, y}] {
		return false
	}

	visited[[2]int{x, y}] = true

	// Recursively check all connected stones
	return hasLiberty(x-1, y, color, visited) ||
		hasLiberty(x+1, y, color, visited) ||
		hasLiberty(x, y-1, color, visited) ||
		hasLiberty(x, y+1, color, visited)
}

// Remove captured stones from the board
func removeCapturedGroup(visited map[[2]int]bool) {
	for pos := range visited {
		board[pos[0]][pos[1]] = "."
	}
}

func main() {
	initBoard()
	fmt.Println("Simple Go Game (9x9) - Now with captures!")
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
