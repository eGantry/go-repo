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
func (g *Game) Update() error {
	// Check for mouse click
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		gridX, gridY := x/cellSize, y/cellSize

		placeStone(gridX, gridY)
		switchPlayer()
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

	// Draw stones as circles
	for i := 0; i < boardSize; i++ {
		for j := 0; j < boardSize; j++ {
			x := float64(i*cellSize) + float64(cellSize)/2
			y := float64(j*cellSize) + float64(cellSize)/2
			radius := float64(cellSize) * 0.4 // Slightly smaller than a grid cell

			if board[i][j] == "B" {
				drawCircle(screen, x, y, radius, color.Black)
			} else if board[i][j] == "W" {
				drawCircle(screen, x, y, radius, color.White)
			}
		}
	}
}

// Draw a filled circle
func drawCircle(screen *ebiten.Image, cx, cy, radius float64, col color.Color) {
	for dy := -radius; dy <= radius; dy++ {
		for dx := -radius; dx <= radius; dx++ {
			if dx*dx+dy*dy <= radius*radius {
				screen.Set(int(cx+dx), int(cy+dy), col)
			}
		}
	}
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

func main() {
	initBoard()
	ebiten.SetWindowSize(windowSize, windowSize)
	ebiten.SetWindowTitle("Go Game with Ebiten")

	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
