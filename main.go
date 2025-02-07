package main

import (
	"fmt"
	"math/rand"
	"time"
)

const size = 9

var board [size][size]string
var currentPlayer = "B" // Player is Black, AI is White (W)

func initBoard() {
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			board[i][j] = "."
		}
	}
}

// Print the board in console
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

// Check if a move is valid
func isValidMove(x, y int) bool {
	if x < 0 || x >= size || y < 0 || y >= size {
		return false
	}
	if board[x][y] != "." {
		return false
	}
	return true
}

// Simple AI: Prioritize captures, otherwise play randomly
func aiMove() {
	time.Sleep(1 * time.Second) // Simulate AI "thinking"

	// 1. Try to find a move that captures an opponent's stone
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			if isValidMove(i, j) {
				// Simulate placing a stone
				board[i][j] = "W"

				// Check if it captures a stone (simplified check)
				if checkCapture(i, j, "W") {
					fmt.Println("AI played at", i, j, "to capture!")
					printBoard()
					currentPlayer = "B"
					return
				}

				// Undo move if it doesn't capture
				board[i][j] = "."
			}
		}
	}

	// 2. Otherwise, play a random valid move
	for {
		x, y := rand.Intn(size), rand.Intn(size)
		if isValidMove(x, y) {
			board[x][y] = "W"
			fmt.Println("AI played at", x, y)
			printBoard()
			currentPlayer = "B"
			return
		}
	}
}

// Simulated capture check (basic version)
func checkCapture(x, y int, color string) bool {
	opponent := "B"
	if color == "B" {
		opponent = "W"
	}

	// Check adjacent spaces for an opponent stone that has no liberties
	directions := [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
	for _, d := range directions {
		nx, ny := x+d[0], y+d[1]
		if nx >= 0 && nx < size && ny >= 0 && ny < size && board[nx][ny] == opponent {
			// If the opponent's stone has no liberties, it's captured
			if !hasLiberty(nx, ny, opponent, make(map[[2]int]bool)) {
				board[nx][ny] = "." // Remove the captured stone
				return true
			}
		}
	}
	return false
}

// Recursively check if a group has liberties
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

func main() {
	rand.Seed(time.Now().UnixNano())
	initBoard()
	fmt.Println("Go Game with Simple AI")
	printBoard()

	for {
		if currentPlayer == "B" {
			var x, y int
			fmt.Print("Your move (row col): ")
			_, err := fmt.Scan(&x, &y)
			if err != nil || !isValidMove(x, y) {
				fmt.Println("Invalid move! Try again.")
				continue
			}
			board[x][y] = "B"
			printBoard()
			currentPlayer = "W"
		} else {
			aiMove()
		}
	}
}
