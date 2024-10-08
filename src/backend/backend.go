package backend

import (
	"math/rand"
)

// TODO: Allow users to provide a custom palette and check if it consists of: min 3 colors; max 8 colors
var DefaultColors = map[string]string{
	"blue":   "#3584e4",
	"green":  "#33d17a",
	"yellow": "#f6d32d",
	"orange": "#ff7800",
	"red":    "#ed333b",
	"purple": "#9141ac",
	//"brown":  "#b5835a",
}

type Position struct {
	Row    int
	Column int
}

type Board struct {
	Rows    int
	Columns int

	Step uint

	Matrix [][]string
}

// InitializeBoard creates a new Board instance with generated "cubicle" matrix of provided size in rows and columns.
func InitializeBoard(rows, cols int) Board {
	matrix := make([][]string, rows)
	for i := range matrix {
		matrix[i] = make([]string, cols)
	}

	var availableColors []string
	for key := range DefaultColors {
		availableColors = append(availableColors, key)
	}

	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			cubicleColor := availableColors[rand.Intn(len(availableColors))]
			matrix[row][col] = cubicleColor
		}
	}

	board := Board{
		Rows:    rows,
		Columns: cols,

		Step: 0,

		Matrix: matrix,
	}

	return board
}

func (b *Board) getNeighbors(pos Position) []Position {
	column := pos.Column
	row := pos.Row

	boardRows := b.Rows
	boardCols := b.Columns

	var neighbors []Position

	// Up
	if row > 0 {
		neighbors = append(neighbors, Position{row - 1, column})
	}

	// Left
	if column > 0 {
		neighbors = append(neighbors, Position{row, column - 1})
	}

	// Down
	if row < (boardRows - 1) {
		neighbors = append(neighbors, Position{row + 1, column})
	}

	// Right
	if column < (boardCols - 1) {
		neighbors = append(neighbors, Position{row, column + 1})
	}

	return neighbors
}

func (b *Board) Flood(newColor string) {
	var queue []Position

	targetColor := b.Matrix[0][0]
	if targetColor == newColor {
		return
	}

	b.Matrix[0][0] = newColor
	queue = append(queue, Position{0, 0})

	for len(queue) > 0 {
		cubicle := queue[0]
		queue = queue[1:]

		neighbor := b.getNeighbors(cubicle)

		for _, pos := range neighbor {
			if b.Matrix[pos.Row][pos.Column] == targetColor {
				b.Matrix[pos.Row][pos.Column] = newColor
				queue = append(queue, pos)
			}
		}
	}

	b.Step++
}

// Formula provided from:
// https://github.com/GunshipPenguin/open_flood/blob/master/app/src/main/java/com/gunshippenguin/openflood/Game.java#L97-L99
func (b *Board) CalculateMaxSteps() uint {
	return uint(30 * (b.Rows * len(DefaultColors)) / (17 * 6))
}

func (b *Board) IsAllFilled() bool {
	targetColor := b.Matrix[0][0]

	for row := 0; row < b.Rows; row++ {
		for col := 0; col < b.Columns; col++ {
			if b.Matrix[row][col] != targetColor {
				return false
			}
		}
	}

	return true
}
