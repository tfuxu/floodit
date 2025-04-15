package views

import (
	"fmt"
	"log/slog"
	"math"

	"github.com/tfuxu/floodit/src/backend/utils"
	"github.com/tfuxu/floodit/src/constants"
	"github.com/tfuxu/floodit/src/views/keyboard"

	"github.com/tfuxu/floodit/src/backend"

	"github.com/diamondburned/gotk4-adwaita/pkg/adw"
	"github.com/diamondburned/gotk4/pkg/cairo"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

type GamePage struct {
	*adw.Bin
	settings *gio.Settings
	parent   *MainWindow

	board backend.Board

	toastOverlay  *adw.ToastOverlay
	gameInfoTitle *adw.WindowTitle

	gameBox     *gtk.Box
	drawingArea *gtk.DrawingArea
}

func NewGamePage(parent *MainWindow, settings *gio.Settings, toastOverlay *adw.ToastOverlay) *GamePage {
	builder := gtk.NewBuilderFromResource(constants.RootPath + "/ui/game_page.ui")

	gamePage := builder.GetObject("game_page").Cast().(*adw.Bin)

	gameInfoTitle := builder.GetObject("game_info_title").Cast().(*adw.WindowTitle)

	gameBox := builder.GetObject("game_box").Cast().(*gtk.Box)
	drawArea := builder.GetObject("draw_area").Cast().(*gtk.DrawingArea)

	board := backend.DefaultBoard()

	gp := GamePage{
		Bin:      gamePage,
		settings: settings,
		parent:   parent,

		board:    board,

		toastOverlay:  toastOverlay,
		gameInfoTitle: gameInfoTitle,

		gameBox:     gameBox,
		drawingArea: drawArea,
	}

	colorKeyboard := keyboard.NewColorKeyboard(backend.DefaultColors, gp.onColorKeyboardUsed)
	gameBox.Append(colorKeyboard)

	/*drawArea.ConnectResize(func(width, height int) {
		fmt.Println("resize", width, height)
	})*/

	drawArea.SetDrawFunc(gp.onDraw)

	return &gp
}

// NewBoard initializes board, sets value for maximum amount of moves
// and queues a board draw.
//
// To get a calculated amount of steps, you need to set the
// `maxSteps` parameter to 0.
//
// To use a random seed, set the `seed` parameter to 0.
func (gp *GamePage) NewBoard(name string, rows, columns int, maxSteps uint, seed int64) {
	gp.board = backend.InitializeBoard(name, rows, columns, seed, maxSteps)

	gp.gameInfoTitle.SetTitle(name)
	// TODO: Translate this
	gp.gameInfoTitle.SetSubtitle(fmt.Sprintf("Steps Left: %d", gp.board.GetStepsLeft()))

	slog.Debug(fmt.Sprintf("maxSteps: %d", gp.board.MaxSteps))
	slog.Debug(fmt.Sprintf("rows: %d columns: %d", gp.board.Rows, gp.board.Columns))

	gp.drawingArea.QueueDraw()
}

func (gp *GamePage) drawBoard(ctx *cairo.Context, width, height int) error {
	boardMatrix := gp.board.Matrix

	boardRows := gp.board.Rows
	boardCols := gp.board.Columns

	rectWidth := width / boardCols
	rectHeight := height / boardRows
	xOffset := (width - rectWidth*boardCols) / 2
	yOffset := (height - rectHeight*boardRows) / 2

	gp.roundedRect(
		ctx,
		float64(xOffset),
		float64(yOffset),
		float64(rectWidth*boardCols),
		float64(rectHeight*boardRows),
		12.0,
	)
	ctx.Clip()

	ctx.NewPath()
	for row := 0; row < boardRows; row++ {
		for col := 0; col < boardCols; col++ {
			x := rectWidth*col + xOffset
			y := rectHeight*row + yOffset
			var hexCode string

			// TODO: This is a dirty workaround to get this working with array.
			// Make sure in future to use array indexes in game matrix instead.
			for _, color := range backend.DefaultColors {
				if color[0] == boardMatrix[row][col] {
					hexCode = color[1]
				}
			}

			cairoRGB, err := utils.HexToCairoRGB(hexCode)
			if err != nil {
				return err
			}

			red := cairoRGB[0]
			green := cairoRGB[1]
			blue := cairoRGB[2]

			ctx.SetSourceRGB(red, green, blue)
			ctx.Rectangle(float64(x), float64(y), float64(rectWidth), float64(rectHeight))
			ctx.Fill()
		}
	}

	return nil
}

func (gp *GamePage) GetCurrentSeed() int64 {
	return gp.board.Seed
}

func (gp *GamePage) onDraw(area *gtk.DrawingArea, ctx *cairo.Context, width, height int) {
	err := gp.drawBoard(ctx, width, height)
	if err != nil {
		gp.toastOverlay.AddToast(adw.NewToast("Failed to retrieve colors for board points"))
		slog.Error("Failed to convert hex values to Cairo compatible RGB channels:", "msg", err)
	}
}

func (gp *GamePage) roundedRect(ctx *cairo.Context, x, y, width, height, cornerRadius float64) {
	ctx.NewSubPath()
	ctx.Arc(x+width-cornerRadius, y+cornerRadius, cornerRadius, -math.Pi/2, 0)
	ctx.Arc(x+width-cornerRadius, y+height-cornerRadius, cornerRadius, 0, math.Pi/2)
	ctx.Arc(x+cornerRadius, y+height-cornerRadius, cornerRadius, math.Pi/2, math.Pi)
	ctx.Arc(x+cornerRadius, y+cornerRadius, cornerRadius, math.Pi, 3*math.Pi/2)
	ctx.ClosePath()
}

func (gp *GamePage) onColorKeyboardUsed(colorName string) {
	if gp.board.GetStepsLeft() < 1 {
		gp.ActivateAction("win.show-finish", glib.NewVariantBoolean(false))
		return
	}

	gp.board.Flood(colorName)

	if gp.board.IsAllFilled() {
		gp.ActivateAction("win.show-finish", glib.NewVariantBoolean(true))
		return
	}

	stepsLeft := gp.board.GetStepsLeft()
	if stepsLeft < 1 {
		gp.ActivateAction("win.show-finish", glib.NewVariantBoolean(false))
		return
	}

	// TODO: Translate this
	gp.gameInfoTitle.SetSubtitle(fmt.Sprintf("Steps Left: %d", stepsLeft))

	slog.Debug(fmt.Sprintf("Step: %d", gp.board.Step))
	slog.Debug(fmt.Sprintf("StepsLeft: %d", stepsLeft))
	slog.Debug(fmt.Sprintf("IsAllFilled: %t", gp.board.IsAllFilled()))

	gp.drawingArea.QueueDraw()
}
