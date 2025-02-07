package main

import (
	"fmt"
	"strings"
)

const size = 9

var board [size][size]string
var currentPlayer = "B" // B = Black, W = White
var previousBoardState string

// Initialize the board
func initBoard() {
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			board[i][j] = "."
		}
	}
	previousBoardState = boardToString()
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

// Convert board state to a string for comparison
func boardToString() string {
	var sb strings.Builder
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			sb.WriteString(board[i][j])
		}
	}
	return sb.String()
}

// Place a stone and check for captures (including Ko rule)
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

	// Check Ko rule: If the new board matches the previous state, undo the move
	newBoardState := boardToString()
	if newBoardState == previousBoardState {
		fmt.Println("Ko rule violated! You cannot immediately recreate the previous board state.")
		board[x][y] = "." // Undo move
		return false
	}

	// Update previous board state for next move
	previousBoardState = newBoardState
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
	fmt.Println("Simple Go Game (9x9) - Now with Ko rule!")
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
