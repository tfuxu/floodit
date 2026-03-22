package views

import (
	"fmt"
	"log/slog"

	"github.com/tfuxu/floodit/src/constants"
	"github.com/tfuxu/floodit/src/views/board"
	"github.com/tfuxu/floodit/src/views/keyboard"

	"github.com/tfuxu/floodit/src/backend"

	"codeberg.org/puregotk/puregotk/v4/adw"
	"codeberg.org/puregotk/puregotk/v4/gio"
	"codeberg.org/puregotk/puregotk/v4/glib"
	"codeberg.org/puregotk/puregotk/v4/gtk"
	. "github.com/pojntfx/go-gettext/pkg/i18n"
)

type GamePage struct {
	*adw.Bin
	settings *gio.Settings
	parent   *MainWindow

	board backend.Board

	toastOverlay  *adw.ToastOverlay
	gameInfoTitle *adw.WindowTitle

	gameBox   *gtk.Box
	gameBoard *board.GameBoard
}

func NewGamePage(parent *MainWindow, settings *gio.Settings, toastOverlay *adw.ToastOverlay) *GamePage {
	builder := gtk.NewBuilderFromResource(constants.RootPath + "/ui/game_page.ui")

	var gamePage adw.Bin
	builder.GetObject("game_page").Cast(&gamePage)
	defer gamePage.Unref()

	var gameInfoTitle adw.WindowTitle
	builder.GetObject("game_info_title").Cast(&gameInfoTitle)
	defer gameInfoTitle.Unref()

	var gameBox gtk.Box
	builder.GetObject("game_box").Cast(&gameBox)
	defer gameBox.Unref()

	defaultBoard := backend.DefaultBoard()

	gp := GamePage{
		Bin:      &gamePage,
		settings: settings,
		parent:   parent,

		board: defaultBoard,

		toastOverlay:  toastOverlay,
		gameInfoTitle: &gameInfoTitle,

		gameBox: &gameBox,
	}

	// TODO: Add a breakpoint that will set a higher content size when window width is >= 600px
	gameBoard := board.NewGameBoard(
		&gp.board,
		"vexpand", true,
		"hexpand", true,
		"width-request", 300,
		"height-request", 300,
	)
	gameBox.Append(&gameBoard.Widget)
	gp.gameBoard = &gameBoard

	colorKeyboard := keyboard.NewColorKeyboard(backend.DefaultColors, gp.onColorKeyboardUsed)
	gameBox.Append(&colorKeyboard.Widget)

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
	// TRANSLATORS: DO NOT translate the '%d' part of the text.
	gp.gameInfoTitle.SetSubtitle(fmt.Sprintf(L("Steps Left: %d"), gp.board.GetStepsLeft()))

	slog.Debug(fmt.Sprintf("maxSteps: %d", gp.board.MaxSteps))
	slog.Debug(fmt.Sprintf("rows: %d columns: %d", gp.board.Rows, gp.board.Columns))

	gp.gameBoard.QueueDraw()
}

func (gp *GamePage) GetCurrentSeed() int64 {
	return gp.board.Seed
}

func (gp *GamePage) onColorKeyboardUsed(colorName string) {
	if gp.board.GetStepsLeft() < 1 {
		gp.ActivateActionVariant("win.show-finish", glib.NewVariantBoolean(false))
		return
	}

	gp.board.Flood(colorName)

	if gp.board.IsAllFilled() {
		gp.ActivateActionVariant("win.show-finish", glib.NewVariantBoolean(true))
		return
	}

	stepsLeft := gp.board.GetStepsLeft()
	if stepsLeft < 1 {
		gp.ActivateActionVariant("win.show-finish", glib.NewVariantBoolean(false))
		return
	}

	// TRANSLATORS: DO NOT translate the '%d' part of the text.
	gp.gameInfoTitle.SetSubtitle(fmt.Sprintf(L("Steps Left: %d"), stepsLeft))

	slog.Debug(fmt.Sprintf("Step: %d", gp.board.Step))
	slog.Debug(fmt.Sprintf("StepsLeft: %d", stepsLeft))
	slog.Debug(fmt.Sprintf("IsAllFilled: %t", gp.board.IsAllFilled()))

	gp.gameBoard.QueueDraw()
}
