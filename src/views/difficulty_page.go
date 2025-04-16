package views

import (
	"github.com/tfuxu/floodit/src/constants"

	"github.com/diamondburned/gotk4-adwaita/pkg/adw"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

type DifficultyPage struct {
	*adw.NavigationPage
	settings *gio.Settings
	parent   *StartingView

	toastOverlay *adw.ToastOverlay

	easyModeRow   *adw.ActionRow
	normalModeRow *adw.ActionRow
	hardModeRow   *adw.ActionRow
	customModeRow *adw.ActionRow
}

func NewDifficultyPage(parent *StartingView, settings *gio.Settings, toastOverlay *adw.ToastOverlay) *DifficultyPage {
	builder := gtk.NewBuilderFromResource(constants.RootPath + "/ui/difficulty_page.ui")

	difficultyPage := builder.GetObject("difficulty_page").Cast().(*adw.NavigationPage)

	easyModeRow := builder.GetObject("easy_mode_row").Cast().(*adw.ActionRow)
	normalModeRow := builder.GetObject("normal_mode_row").Cast().(*adw.ActionRow)
	hardModeRow := builder.GetObject("hard_mode_row").Cast().(*adw.ActionRow)
	customModeRow := builder.GetObject("custom_mode_row").Cast().(*adw.ActionRow)

	cmp := DifficultyPage{
		NavigationPage: difficultyPage,
		settings:       settings,
		parent:         parent,

		toastOverlay: toastOverlay,

		easyModeRow:   easyModeRow,
		normalModeRow: normalModeRow,
		hardModeRow:   hardModeRow,
		customModeRow: customModeRow,
	}

	cmp.setupSignals()

	return &cmp
}

func (dp *DifficultyPage) setupSignals() {
	dp.easyModeRow.ConnectActivated(dp.onEasyModeRowActivated)
	dp.normalModeRow.ConnectActivated(dp.onNormalModeRowActivated)
	dp.hardModeRow.ConnectActivated(dp.onHardModeRowActivated)
	dp.customModeRow.ConnectActivated(dp.onCustomModeRowActivated)
}

/*func (dp *DifficultyPage) onBackButtonClicked() {
	dp.parent.parent.ShowWelcomePage()

}*/

func (dp *DifficultyPage) onEasyModeRowActivated() {
	dp.parent.parent.StartNewGame("Easy", 6, 6, 15, 0) // TODO: Make this translatable
}

func (dp *DifficultyPage) onNormalModeRowActivated() {
	dp.parent.parent.StartNewGame("Normal", 10, 10, 20, 0) // TODO: Make this translatable
}

func (dp *DifficultyPage) onHardModeRowActivated() {
	dp.parent.parent.StartNewGame("Hard", 14, 14, 25, 0) // TODO: Make this translatable
}

func (dp *DifficultyPage) onCustomModeRowActivated() {
	dp.parent.PushByTag("custom-mode")
}
