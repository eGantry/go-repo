package main

var previousBoardState string

// **Converts board to a string (for Ko rule)**
func BoardToString() string {
	var str string
	for i := 0; i < boardSize; i++ {
		for j := 0; j < boardSize; j++ {
			str += board[i][j]
		}
	}
	return str
}

// **Enforce the Ko rule**
func EnforceKoRule() bool {
	tempBoard := BoardToString()
	if tempBoard == previousBoardState {
		return false // Move would recreate previous state
	}
	previousBoardState = tempBoard
	return true
}

// **Check for captured stones**
func CheckForCaptures() {
	opponent := "B"
	if currentPlayer == "B" {
		opponent = "W"
	}

	for i := 0; i < boardSize; i++ {
		for j := 0; j < boardSize; j++ {
			if board[i][j] == opponent {
				visited := make(map[[2]int]bool)
				if !HasLiberty(i, j, opponent, visited) {
					RemoveCapturedGroup(visited)
				}
			}
		}
	}
}

// **Check for liberties (empty spaces)**
func HasLiberty(x, y int, color string, visited map[[2]int]bool) bool {
	if x < 0 || x >= boardSize || y < 0 || y >= boardSize {
		return false
	}
	if board[x][y] == "." {
		return true
	}
	if board[x][y] != color || visited[[2]int{x, y}] {
		return false
	}

	visited[[2]int{x, y}] = true

	return HasLiberty(x-1, y, color, visited) ||
		HasLiberty(x+1, y, color, visited) ||
		HasLiberty(x, y-1, color, visited) ||
		HasLiberty(x, y+1, color, visited)
}

// **Remove captured stones**
func RemoveCapturedGroup(visited map[[2]int]bool) {
	for pos := range visited {
		board[pos[0]][pos[1]] = "."
	}
}
