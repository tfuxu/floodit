package views

import (
	"fmt"
	"log/slog"

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

	board     backend.Board
	maxSteps  uint

	toastOverlay  *adw.ToastOverlay
	gameInfoTitle *adw.WindowTitle

	gameBox     *gtk.Box
	boardView   *gtk.Box
	drawingArea *gtk.DrawingArea
}

type BoardConfig struct {
	Rows    int
	Columns int

	MaxSteps int
}

func NewGamePage(parent *MainWindow, settings *gio.Settings, toastOverlay *adw.ToastOverlay) *GamePage {
	builder := gtk.NewBuilderFromResource(constants.RootPath + "/ui/game_page.ui")

	gamePage := builder.GetObject("game_page").Cast().(*adw.Bin)

	gameInfoTitle := builder.GetObject("game_info_title").Cast().(*adw.WindowTitle)

	gameBox := builder.GetObject("game_box").Cast().(*gtk.Box)
	boardView := builder.GetObject("board").Cast().(*gtk.Box)
	drawArea := builder.GetObject("draw_area").Cast().(*gtk.DrawingArea)

	board := backend.InitializeBoard(10, 10)
	maxSteps := board.CalculateMaxSteps()

	gp := GamePage{
		Bin:      gamePage,
		settings: settings,
		parent:   parent,

		board:     board,
		maxSteps:  maxSteps,

		toastOverlay:  toastOverlay,
		gameInfoTitle: gameInfoTitle,

		gameBox:     gameBox,
		boardView:   boardView,
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

// NewBoard initializes board, sets value for max amount of moves and queues a board draw
// NOTE: To get calculated amount of steps, you need to set maxSteps parameter to 0.
func (gp *GamePage) NewBoard(name string, rows int, cols int, maxSteps uint) {
	gp.board = backend.InitializeBoard(rows, cols)

	if maxSteps == 0 {
		gp.maxSteps = gp.board.CalculateMaxSteps()
	} else {
		gp.maxSteps = maxSteps
	}

	gp.gameInfoTitle.SetTitle(name)
	// TODO: Translate this
	gp.gameInfoTitle.SetSubtitle(fmt.Sprintf("Steps Left: %d", gp.maxSteps))

	slog.Debug(fmt.Sprintf("maxSteps: %d", gp.maxSteps))
	slog.Debug(fmt.Sprintf("rows: %d columns: %d", gp.board.Rows, gp.board.Columns))

	gp.drawingArea.QueueDraw()
}

func (gp *GamePage) drawBoard(ctx *cairo.Context, width, height int) error {
	boardMatrix := gp.board.Matrix

	boardRows := gp.board.Rows
	boardCols := gp.board.Columns

	rectWidth := float64(width) / float64(boardCols)
	rectHeight := float64(height) / float64(boardRows)

	for row := 0; row < boardRows; row++ {
		for col := 0; col < boardCols; col++ {
			x := float64(rectWidth) * float64(col)
			y := float64(rectHeight) * float64(row)

			hexCode := backend.DefaultColors[boardMatrix[row][col]]
			cairoRGB, err := utils.HexToCairoRGB(hexCode)
			if err != nil {
				return err
			}

			red := cairoRGB[0]
			green := cairoRGB[1]
			blue := cairoRGB[2]

			ctx.SetSourceRGB(red, green, blue)
			ctx.Rectangle(x, y, float64(rectWidth), float64(rectHeight))
			ctx.Fill()
		}
	}

	return nil
}

func (gp *GamePage) onDraw(area *gtk.DrawingArea, ctx *cairo.Context, width, height int) {
	err := gp.drawBoard(ctx, width, height)
	if err != nil {
		gp.toastOverlay.AddToast(adw.NewToast("Failed to retrieve colors for board points"))
		slog.Error("Failed to convert hex values to Cairo compatible RGB channels:", "msg", err)
	}
}

func (gp *GamePage) onColorKeyboardUsed(colorName string) {
	gp.board.Flood(colorName)

	if gp.board.IsAllFilled() {
		gp.ActivateAction("win.show-finish", glib.NewVariantBoolean(true))
		return
	}

	stepsLeft := int(gp.maxSteps - gp.board.Step)
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
