package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	boardSize  = 9
	cellSize   = 50
	windowSize = boardSize * cellSize
)

var board [boardSize][boardSize]string
var currentPlayer = "B" // B = Black, W = White

var passCount int // Track consecutive passes

// Handle passing turns
func passTurn() {
	passCount++
	if passCount >= 2 {
		fmt.Println("Both players passed. Ending game...")
		calculateFinalScore()
		return
	}
	switchPlayer()
}

func initBoard() {
	for i := 0; i < boardSize; i++ {
		for j := 0; j < boardSize; j++ {
			board[i][j] = "."
		}
	}
}

// **Ebiten Game Structure**
type Game struct{}

// **Update handles input**
var isPassing bool // Track if a player is attempting to pass

func (g *Game) Update() error {
	// If player is confirming a pass
	if isPassing {
		if ebiten.IsKeyPressed(ebiten.KeyY) {
			fmt.Println("Player", currentPlayer, "confirmed pass!")
			passTurn()
			isPassing = false
		} else if ebiten.IsKeyPressed(ebiten.KeyN) {
			fmt.Println("Pass canceled!")
			isPassing = false
		}
		return nil
	}

	// Handle pass request (first press)
	if ebiten.IsKeyPressed(ebiten.KeyP) {
		fmt.Println("Player", currentPlayer, "wants to pass. Press 'Y' to confirm or 'N' to cancel.")
		isPassing = true
		return nil
	}

	// Handle mouse click for placing stones
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		gridX, gridY := x/cellSize, y/cellSize

		if isValidMove(gridX, gridY) {
			board[gridX][gridY] = currentPlayer
			passCount = 0 // Reset pass count if a move is made
			switchPlayer()
		}
	}
	return nil
}

func placeStone(x, y int) bool {
	if !isValidMove(x, y) {
		return false
	}

	// Place the stone
	board[x][y] = currentPlayer

	// Enforce the Ko rule: If the move recreates the previous state, undo it.
	if !EnforceKoRule() {
		fmt.Println("Ko rule violated! Move reverted.")
		board[x][y] = "." // Undo the move
		return false
	}

	// Check for captures on both sides
	CheckForCaptures()

	return true
}

// Draw renders the board and stones
func (g *Game) Draw(screen *ebiten.Image) {
	// Set background color (light wood)
	screen.Fill(color.RGBA{R: 222, G: 184, B: 135, A: 255}) // Light brown board

	// Draw grid
	for i := 0; i < boardSize; i++ {
		for j := 0; j < boardSize; j++ {
			ebitenutil.DrawLine(screen, float64(i*cellSize), 0, float64(i*cellSize), windowSize, color.Black)
			ebitenutil.DrawLine(screen, 0, float64(i*cellSize), windowSize, float64(i*cellSize), color.Black)
		}
	}

	// Draw stones as 3D shaded circles
	for i := 0; i < boardSize; i++ {
		for j := 0; j < boardSize; j++ {
			x := float64(i*cellSize) + float64(cellSize)/2
			y := float64(j*cellSize) + float64(cellSize)/2
			radius := float64(cellSize) * 0.4 // Slightly smaller than a grid cell

			if board[i][j] == "B" {
				drawStone(screen, x, y, radius, color.Black)
			} else if board[i][j] == "W" {
				drawStone(screen, x, y, radius, color.White)
			}
		}
	}

}

// Draw a shaded, 3D-style stone
func drawStone(screen *ebiten.Image, cx, cy, radius float64, col color.Color) {
	for dy := -radius; dy <= radius; dy++ {
		for dx := -radius; dx <= radius; dx++ {
			dist := dx*dx + dy*dy
			if dist <= radius*radius {
				// Create a shading effect: lighter at the top-left, darker at bottom-right
				shadeFactor := (dx + dy) / radius * 0.3
				baseColor := col

				// Convert base color to RGBA for manipulation
				r, g, b, a := colorToRGBA(baseColor)
				r = clamp(r + int(50*shadeFactor))
				g = clamp(g + int(50*shadeFactor))
				b = clamp(b + int(50*shadeFactor))

				// Apply highlight near the top-left
				if dist < (radius*radius)*0.3 {
					r = clamp(r + 60)
					g = clamp(g + 60)
					b = clamp(b + 60)
				}

				// Apply the pixel with shading
				screen.Set(int(cx+dx), int(cy+dy), color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)})
			}
		}
	}
}

// Helper function: Convert color to RGBA components
func colorToRGBA(c color.Color) (int, int, int, int) {
	r, g, b, a := c.RGBA()
	return int(r >> 8), int(g >> 8), int(b >> 8), int(a >> 8)
}

// Helper function: Ensure RGB values stay in valid range
func clamp(value int) int {
	if value > 255 {
		return 255
	} else if value < 0 {
		return 0
	}
	return value
}

// **Layout defines screen size**
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return windowSize, windowSize
}

// **Check for a valid move**
func isValidMove(x, y int) bool {
	if x < 0 || x >= boardSize || y < 0 || y >= boardSize {
		return false
	}
	return board[x][y] == "."
}

// **Switch player after a move**
func switchPlayer() {
	if currentPlayer == "B" {
		currentPlayer = "W"
	} else {
		currentPlayer = "B"
	}
}

func calculateFinalScore() {
	blackScore, whiteScore := calculateTerritory()
	fmt.Printf("\nFinal Scores:\nBlack: %d\nWhite: %d\n", blackScore, whiteScore)

	if blackScore > whiteScore {
		fmt.Println("🏆 Black wins!")
	} else if whiteScore > blackScore {
		fmt.Println("🏆 White wins!")
	} else {
		fmt.Println("🤝 It's a tie!")
	}

	// Exit game after scoring
	ebiten.SetRunnableOnUnfocused(false) // Ensure window closes cleanly
	log.Fatal("Game Over!")
}

// Count controlled territory
func calculateTerritory() (int, int) {
	visited := make(map[[2]int]bool)
	blackScore, whiteScore := 0, 0

	for i := 0; i < boardSize; i++ {
		for j := 0; j < boardSize; j++ {
			if board[i][j] == "." && !visited[[2]int{i, j}] {
				territory, owner := identifyTerritory(i, j, visited)
				if owner == "B" {
					blackScore += territory
				} else if owner == "W" {
					whiteScore += territory
				}
			}
		}
	}
	return blackScore, whiteScore
}

// Identify empty spaces and determine which player surrounds them
func identifyTerritory(x, y int, visited map[[2]int]bool) (int, string) {
	stack := [][2]int{{x, y}}
	territory := 0
	owner := ""

	for len(stack) > 0 {
		pos := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		i, j := pos[0], pos[1]
		if i < 0 || i >= boardSize || j < 0 || j >= boardSize || visited[[2]int{i, j}] {
			continue
		}
		if board[i][j] == "B" {
			if owner == "W" {
				return 0, ""
			}
			owner = "B"
			continue
		}
		if board[i][j] == "W" {
			if owner == "B" {
				return 0, ""
			}
			owner = "W"
			continue
		}

		territory++
		visited[[2]int{i, j}] = true
		stack = append(stack, [2]int{i - 1, j}, [2]int{i + 1, j}, [2]int{i, j - 1}, [2]int{i, j + 1})
	}

	return territory, owner
}

func main() {
	initBoard()
	ebiten.SetWindowSize(windowSize, windowSize)
	ebiten.SetWindowTitle("Go Game with Ebiten")

	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
