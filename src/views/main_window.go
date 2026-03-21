package views

import (
	"unsafe"

	"github.com/tfuxu/floodit/src/constants"

	"codeberg.org/puregotk/puregotk/v4/adw"
	"codeberg.org/puregotk/puregotk/v4/gio"
	"codeberg.org/puregotk/puregotk/v4/glib"
	"codeberg.org/puregotk/puregotk/v4/gtk"
)

type MainWindow struct {
	*adw.ApplicationWindow
	app      *adw.Application
	settings *gio.Settings

	toastOverlay *adw.ToastOverlay
	mainStack    *gtk.Stack

	statusPage *adw.StatusPage
	playButton *gtk.Button
	//guideButton *gtk.Button

	startingView  *StartingView
	gamePage      *GamePage
	gameRulesPage *GameRulesPage
	resultPage    *ResultPage
}

func NewMainWindow(app *adw.Application, settings *gio.Settings) *MainWindow {
	builder := gtk.NewBuilderFromResource(constants.RootPath + "/ui/main_window.ui")

	var window adw.ApplicationWindow
	builder.GetObject("main_window").Cast(&window)
	defer window.Unref()

	window.SetApplication(&app.Application)
	window.SetDefaultSize(
		settings.GetInt("window-width"), settings.GetInt("window-height"),
	)

	var statusPage adw.StatusPage
	builder.GetObject("status_page").Cast(&statusPage)
	defer statusPage.Unref()

	var playButton gtk.Button
	builder.GetObject("play_button").Cast(&playButton)
	defer playButton.Unref()

	//var guideButton gtk.Button
	//builder.GetObject("guide_button").Cast(&guideButton)
	//defer guideButton.Unref()

	var mainStack gtk.Stack
	builder.GetObject("main_stack").Cast(&mainStack)
	defer mainStack.Unref()

	var toastOverlay adw.ToastOverlay
	builder.GetObject("toast_overlay").Cast(&toastOverlay)
	defer toastOverlay.Unref()

	w := MainWindow{
		ApplicationWindow: &window,
		app:               app,
		settings:          settings,

		statusPage: &statusPage,
		playButton: &playButton,
		//guideButton: &guideButton,

		toastOverlay: &toastOverlay,
		mainStack:    &mainStack,
	}
	w.startingView = NewStartingView(&w, settings, &toastOverlay)
	w.gamePage = NewGamePage(&w, settings, &toastOverlay)
	w.gameRulesPage = NewGameRulesPage(&w, settings, &toastOverlay)
	w.resultPage = NewResultPage(&w, settings, &toastOverlay)

	statusPage.SetIconName(constants.AppID)

	if constants.BuildType == "debug" {
		w.AddCssClass("devel")
	}

	w.setupActions()
	w.setupSignals()
	w.setupStack()

	return &w
}

func (w *MainWindow) setupActions() {
	// TODO: Disable this action in certain parts of the app to prevent undefined behavior
	playAgainAction := gio.NewSimpleAction("play-again", nil)
	playAgainAction.ConnectActivate(new(func(gio.SimpleAction, uintptr) {
		w.playAgain()
	}))
	w.AddAction(playAgainAction)

	showWelcomeAction := gio.NewSimpleAction("show-welcome", nil)
	showWelcomeAction.ConnectActivate(new(func(gio.SimpleAction, uintptr) {
		w.showWelcomePage()
	}))
	w.AddAction(showWelcomeAction)

	showGameSelectAction := gio.NewSimpleAction("show-game-select", nil)
	showGameSelectAction.ConnectActivate(new(func(gio.SimpleAction, uintptr) {
		w.startingView.Pop()
		w.showStartingPage()
	}))
	w.AddAction(showGameSelectAction)

	showGameRulesAction := gio.NewSimpleAction("show-game-rules", nil)
	showGameRulesAction.ConnectActivate(new(func(gio.SimpleAction, uintptr) {
		w.showRulesPage()
	}))
	w.AddAction(showGameRulesAction)

	showFinishAction := gio.NewSimpleAction("show-finish", glib.NewVariantType("b"))
	showFinishAction.ConnectActivate(new(func(_ gio.SimpleAction, parameter uintptr) {
		variant := (*glib.Variant)(unsafe.Pointer(parameter))
		w.showResultPage(variant.GetBoolean())
	}))
	w.AddAction(showFinishAction)
}

func (w *MainWindow) setupSignals() {
	w.playButton.ConnectClicked(new(func(gtk.Button) {
		w.onPlayClicked()
	}))

	w.ConnectUnrealize(new(func(gtk.Widget) {
		w.saveWindowProps()
	}))
}

func (w *MainWindow) setupStack() {
	w.mainStack.AddNamed(&w.startingView.Widget, "stack_starting_page")
	w.mainStack.AddNamed(&w.gamePage.Widget, "stack_game_page")
	w.mainStack.AddNamed(&w.gameRulesPage.Widget, "stack_game_rules_page")
	w.mainStack.AddNamed(&w.resultPage.Widget, "stack_result_page")
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
	var width, height int32
	w.GetDefaultSize(&width, &height)

	w.settings.SetInt("window-width", width)
	w.settings.SetInt("window-height", height)
	w.settings.SetBoolean("window-maximized", w.IsMaximized())
}
