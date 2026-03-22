package views

import (
	"github.com/tfuxu/floodit/src/constants"

	"codeberg.org/puregotk/puregotk/v4/adw"
	"codeberg.org/puregotk/puregotk/v4/gio"
	"codeberg.org/puregotk/puregotk/v4/gtk"
	. "github.com/pojntfx/go-gettext/pkg/i18n"
)

type CustomModePage struct {
	*adw.NavigationPage
	settings *gio.Settings
	parent   *StartingView

	toastOverlay *adw.ToastOverlay

	boardSizeRow             *adw.SpinRow
	enableCustomMoveLimitRow *adw.SwitchRow
	moveLimitRow             *adw.SpinRow

	playButton *gtk.Button
}

func NewCustomModePage(parent *StartingView, settings *gio.Settings, toastOverlay *adw.ToastOverlay) *CustomModePage {
	builder := gtk.NewBuilderFromResource(constants.RootPath + "/ui/custom_mode_page.ui")

	var customModePage adw.NavigationPage
	builder.GetObject("custom_mode_page").Cast(&customModePage)
	defer customModePage.Unref()

	var boardSizeRow adw.SpinRow
	builder.GetObject("board_size_row").Cast(&boardSizeRow)
	defer boardSizeRow.Unref()

	var enableCustomMoveLimitRow adw.SwitchRow
	builder.GetObject("enable_custom_move_limit_row").Cast(&enableCustomMoveLimitRow)
	defer enableCustomMoveLimitRow.Unref()

	var moveLimitRow adw.SpinRow
	builder.GetObject("move_limit_row").Cast(&moveLimitRow)
	defer moveLimitRow.Unref()

	var playButton gtk.Button
	builder.GetObject("play_button").Cast(&playButton)
	defer playButton.Unref()

	cmp := CustomModePage{
		NavigationPage: &customModePage,
		settings:       settings,
		parent:         parent,

		toastOverlay: toastOverlay,

		boardSizeRow:             &boardSizeRow,
		enableCustomMoveLimitRow: &enableCustomMoveLimitRow,
		moveLimitRow:             &moveLimitRow,

		playButton: &playButton,
	}

	// Workaround: Set default values for SpinRows
	cmp.boardSizeRow.SetValue(2)
	cmp.moveLimitRow.SetValue(1)

	cmp.setupSignals()

	return &cmp
}

func (cmp *CustomModePage) setupSignals() {
	cmp.enableCustomMoveLimitRow.ConnectSignal("notify::active", new(func() {
		cmp.onEnableCustomMoveLimitSwitched()
	}))

	cmp.playButton.ConnectClicked(new(func(gtk.Button) {
		cmp.onPlayButtonClicked()
	}))
}

func (cmp *CustomModePage) onEnableCustomMoveLimitSwitched() {
	if cmp.enableCustomMoveLimitRow.GetActive() {
		cmp.moveLimitRow.SetSensitive(true)
	} else {
		cmp.moveLimitRow.SetSensitive(false)
	}
}

func (cmp *CustomModePage) onPlayButtonClicked() {
	rows := int(cmp.boardSizeRow.GetValue())
	cols := int(cmp.boardSizeRow.GetValue())

	var maxSteps uint

	if cmp.enableCustomMoveLimitRow.GetActive() {
		maxSteps = uint(cmp.moveLimitRow.GetValue())
	} else {
		maxSteps = 0
	}

	cmp.parent.parent.StartNewGame(L("Custom"), rows, cols, maxSteps, 0)
}
