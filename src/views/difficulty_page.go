package views

import (
	"github.com/tfuxu/floodit/src/constants"

	"codeberg.org/puregotk/puregotk/v4/adw"
	"codeberg.org/puregotk/puregotk/v4/gio"
	"codeberg.org/puregotk/puregotk/v4/gtk"
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

	var difficultyPage adw.NavigationPage
	builder.GetObject("difficulty_page").Cast(&difficultyPage)
	defer difficultyPage.Unref()

	var easyModeRow adw.ActionRow
	builder.GetObject("easy_mode_row").Cast(&easyModeRow)
	defer easyModeRow.Unref()

	var normalModeRow adw.ActionRow
	builder.GetObject("normal_mode_row").Cast(&normalModeRow)
	defer normalModeRow.Unref()

	var hardModeRow adw.ActionRow
	builder.GetObject("hard_mode_row").Cast(&hardModeRow)
	defer hardModeRow.Unref()

	var customModeRow adw.ActionRow
	builder.GetObject("custom_mode_row").Cast(&customModeRow)
	defer customModeRow.Unref()

	cmp := DifficultyPage{
		NavigationPage: &difficultyPage,
		settings:       settings,
		parent:         parent,

		toastOverlay: toastOverlay,

		easyModeRow:   &easyModeRow,
		normalModeRow: &normalModeRow,
		hardModeRow:   &hardModeRow,
		customModeRow: &customModeRow,
	}

	cmp.setupSignals()

	return &cmp
}

func (dp *DifficultyPage) setupSignals() {
	dp.easyModeRow.ConnectActivated(new(func(adw.ActionRow) { dp.onEasyModeRowActivated() }))
	dp.normalModeRow.ConnectActivated(new(func(adw.ActionRow) { dp.onNormalModeRowActivated() }))
	dp.hardModeRow.ConnectActivated(new(func(adw.ActionRow) { dp.onHardModeRowActivated() }))
	dp.customModeRow.ConnectActivated(new(func(adw.ActionRow) { dp.onCustomModeRowActivated() }))
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
