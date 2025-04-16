package backend

import (
	"math/rand"
	"time"
)

// TODO: Allow users to provide a custom palette and check if it consists of:
// min 3 colors; max 8 colors
var DefaultColors = [][2]string{
	{"blue",   "#3584e4"},
	{"green",  "#33d17a"},
	{"yellow", "#f6d32d"},
	{"orange", "#ff7800"},
	{"red",    "#ed333b"},
	{"purple", "#9141ac"},
	//{"brown",  "#b5835a"},
}

type Position struct {
	Row    int
	Column int
}

type Board struct {
	Name string
	Seed int64

	Rows    int
	Columns int

	Step     uint
	MaxSteps uint

	// TODO: Make matrix take color array indexes instead
	Matrix [][]string
}

// Creates a new empty Board instance to use when initializing stuff
//
// NOTE: Board matrix is set to nil, so make sure to fill it with data
// before making operations on it.
func DefaultBoard() Board {
	b := Board{
		Name: "Custom",

		Rows: 0,
		Columns: 0,

		Step: 0,
		MaxSteps: 1,
	}

	return b
}

// InitializeBoard creates a new Board instance with generated "cubicle"
// matrix of provided size in rows and columns.
//
// To get a calculated amount of steps, you need to set the
// `maxSteps` parameter to 0.
//
// To use a random seed, set the `seed` parameter to 0.
func InitializeBoard(name string, rows, columns int, seed int64, maxSteps uint) Board {
	matrix := make([][]string, rows)
	for i := range matrix {
		matrix[i] = make([]string, columns)
	}

	var availableColors []string
	for _, value := range DefaultColors {
		availableColors = append(availableColors, value[0])
	}

	if seed == 0 {
		seed = time.Now().UnixNano()
	}

	random := rand.New(rand.NewSource(seed))

	for row := 0; row < rows; row++ {
		for col := 0; col < columns; col++ {
			cubicleColor := availableColors[random.Intn(len(availableColors))]
			matrix[row][col] = cubicleColor
		}
	}

	board := Board{
		Name: name,
		Seed: seed,

		Rows:    rows,
		Columns: columns,

		Step: 0,

		Matrix: matrix,
	}

	if maxSteps == 0 {
		maxSteps = board.CalculateMaxSteps()
	}

	board.MaxSteps = maxSteps

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

// TODO: Check if this shouldn't use uint instead
func (b *Board) GetStepsLeft() int {
	return int(b.MaxSteps - b.Step)
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
