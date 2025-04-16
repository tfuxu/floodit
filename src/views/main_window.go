package views

import (
	"github.com/tfuxu/floodit/src/constants"

	"github.com/diamondburned/gotk4-adwaita/pkg/adw"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

type MainWindow struct {
	*adw.ApplicationWindow
	app      *adw.Application
	settings *gio.Settings

	toastOverlay *adw.ToastOverlay
	mainStack    *gtk.Stack

	statusPage *adw.StatusPage
	playButton *gtk.Button

	startingView  *StartingView
	gamePage      *GamePage
	gameRulesPage *GameRulesPage
	resultPage    *ResultPage
}

func NewMainWindow(app *adw.Application, settings *gio.Settings) *MainWindow {
	builder := gtk.NewBuilderFromResource(constants.RootPath + "/ui/main_window.ui")

	window := builder.GetObject("main_window").Cast().(*adw.ApplicationWindow)
	window.SetApplication(&app.Application)
	window.SetDefaultSize(
		settings.Int("window-width"), settings.Int("window-height"),
	)

	statusPage := builder.GetObject("status_page").Cast().(*adw.StatusPage)
	playButton := builder.GetObject("play_button").Cast().(*gtk.Button)

	mainStack := builder.GetObject("main_stack").Cast().(*gtk.Stack)
	toastOverlay := builder.GetObject("toast_overlay").Cast().(*adw.ToastOverlay)

	w := MainWindow{
		ApplicationWindow: window,
		app:               app,
		settings:          settings,

		statusPage: statusPage,
		playButton: playButton,

		toastOverlay: toastOverlay,
		mainStack:    mainStack,
	}
	w.startingView = NewStartingView(&w, settings, toastOverlay)
	w.gamePage = NewGamePage(&w, settings, toastOverlay)
	w.gameRulesPage = NewGameRulesPage(&w, settings, toastOverlay)
	w.resultPage = NewResultPage(&w, settings, toastOverlay)

	statusPage.SetIconName(constants.AppID)

	if constants.BuildType == "debug" {
		w.AddCSSClass("devel")
	}

	w.setupActions()
	w.setupSignals()
	w.setupStack()

	return &w
}

func (w *MainWindow) setupActions() {
	// TODO: Disable this action in certain parts of the app to prevent undefined behavior
	playAgainAction := gio.NewSimpleAction("play-again", nil)
	playAgainAction.ConnectActivate(func(parameter *glib.Variant) {
		w.playAgain()
	})
	w.AddAction(playAgainAction)

	showWelcomeAction := gio.NewSimpleAction("show-welcome", nil)
	showWelcomeAction.ConnectActivate(func(parameter *glib.Variant) {
		w.showWelcomePage()
	})
	w.AddAction(showWelcomeAction)

	showGameSelectAction := gio.NewSimpleAction("show-game-select", nil)
	showGameSelectAction.ConnectActivate(func(parameter *glib.Variant) {
		w.startingView.Pop()
		w.showStartingPage()
	})
	w.AddAction(showGameSelectAction)

	showGameRulesAction := gio.NewSimpleAction("show-game-rules", nil)
	showGameRulesAction.ConnectActivate(func(parameter *glib.Variant) {
		w.showRulesPage()
	})
	w.AddAction(showGameRulesAction)

	showFinishAction := gio.NewSimpleAction("show-finish", glib.NewVariantType("b"))
	showFinishAction.ConnectActivate(func(parameter *glib.Variant) {
		w.showResultPage(parameter.Boolean())
	})
	w.AddAction(showFinishAction)
}

func (w *MainWindow) setupSignals() {
	w.playButton.ConnectClicked(
		w.onPlayClicked,
	)

	w.ConnectUnrealize(func() {
		w.saveWindowProps()
	})
}

func (w *MainWindow) setupStack() {
	w.mainStack.AddNamed(w.startingView, "stack_starting_page")
	w.mainStack.AddNamed(w.gamePage, "stack_game_page")
	w.mainStack.AddNamed(w.gameRulesPage, "stack_game_rules_page")
	w.mainStack.AddNamed(w.resultPage, "stack_result_page")
}

// StartNewGame calls NewBoard to initialize board, and changes view in
// main stack to show game page.
//
// To get a calculated amount of steps, you need to set the
// `maxSteps` parameter to 0.
//
// To use a random seed, set the `seed` parameter to 0.
func (w *MainWindow) StartNewGame(name string, rows, columns int, maxSteps uint, seed int64) {
	w.gamePage.NewBoard(name, rows, columns, maxSteps, seed)
	w.mainStack.SetVisibleChildName("stack_game_page")
}

func (w *MainWindow) playAgain() {
	w.StartNewGame(
		w.gamePage.board.Name,
		w.gamePage.board.Rows,
		w.gamePage.board.Columns,
		w.gamePage.board.MaxSteps,
		w.gamePage.board.Seed,
	)
}

func (w *MainWindow) showWelcomePage() {
	w.mainStack.SetVisibleChildName("stack_welcome_page")
}

func (w *MainWindow) showStartingPage() {
	w.mainStack.SetVisibleChildName("stack_starting_page")
}

func (w *MainWindow) showRulesPage() {
	w.mainStack.SetVisibleChildName("stack_game_rules_page")
}

func (w *MainWindow) showResultPage(isWin bool) {
	w.resultPage.SetResultState(isWin)
	w.mainStack.SetVisibleChildName("stack_result_page")
}

func (w *MainWindow) onPlayClicked() {
	w.showStartingPage()
}

func (w *MainWindow) saveWindowProps() {
	width, height := w.DefaultSize()

	w.settings.SetInt("window-width", width)
	w.settings.SetInt("window-height", height)
	w.settings.SetBoolean("window-maximized", w.IsMaximized())
}
