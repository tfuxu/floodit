package views

import (
	"github.com/tfuxu/floodit/src/constants"

	"github.com/diamondburned/gotk4-adwaita/pkg/adw"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
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

	customModePage := builder.GetObject("custom_mode_page").Cast().(*adw.NavigationPage)

	boardSizeRow := builder.GetObject("board_size_row").Cast().(*adw.SpinRow)
	enableCustomMoveLimitRow := builder.GetObject("enable_custom_move_limit_row").Cast().(*adw.SwitchRow)
	moveLimitRow := builder.GetObject("move_limit_row").Cast().(*adw.SpinRow)

	playButton := builder.GetObject("play_button").Cast().(*gtk.Button)

	cmp := CustomModePage{
		NavigationPage: customModePage,
		settings:       settings,
		parent:         parent,

		toastOverlay: toastOverlay,

		boardSizeRow:             boardSizeRow,
		enableCustomMoveLimitRow: enableCustomMoveLimitRow,
		moveLimitRow:             moveLimitRow,

		playButton: playButton,
	}

	// Workaround: Set default values for SpinRows
	cmp.boardSizeRow.SetValue(2)
	cmp.moveLimitRow.SetValue(1)

	cmp.setupSignals()

	return &cmp
}

func (cmp *CustomModePage) setupSignals() {
	cmp.enableCustomMoveLimitRow.Connect("notify::active", func() {
		cmp.onEnableCustomMoveLimitSwitched()
	})

	cmp.playButton.ConnectClicked(func() {
		cmp.onPlayButtonClicked()
	})
}

func (cmp *CustomModePage) onEnableCustomMoveLimitSwitched() {
	if cmp.enableCustomMoveLimitRow.Active() {
		cmp.moveLimitRow.SetSensitive(true)
	} else {
		cmp.moveLimitRow.SetSensitive(false)
	}
}

func (cmp *CustomModePage) onPlayButtonClicked() {
	rows := int(cmp.boardSizeRow.Value())
	cols := int(cmp.boardSizeRow.Value())

	var maxSteps uint

	if cmp.enableCustomMoveLimitRow.Active() {
		maxSteps = uint(cmp.moveLimitRow.Value())
	} else {
		maxSteps = 0
	}

	cmp.parent.parent.StartNewGame("Custom", rows, cols, maxSteps) // TODO: Make this translatable
}
