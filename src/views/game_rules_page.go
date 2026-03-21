package views

import (
	"github.com/tfuxu/floodit/src/constants"

	"codeberg.org/puregotk/puregotk/v4/adw"
	"codeberg.org/puregotk/puregotk/v4/gio"
	"codeberg.org/puregotk/puregotk/v4/gtk"
)

type GameRulesPage struct {
	*adw.NavigationPage
	settings *gio.Settings
	parent   *MainWindow

	toastOverlay *adw.ToastOverlay
}

func NewGameRulesPage(parent *MainWindow, settings *gio.Settings, toastOverlay *adw.ToastOverlay) *GameRulesPage {
	builder := gtk.NewBuilderFromResource(constants.RootPath + "/ui/game_rules_page.ui")

	var gameGuidePage adw.NavigationPage
	builder.GetObject("game_guide_page").Cast(&gameGuidePage)
	defer gameGuidePage.Unref()

	grp := GameRulesPage{
		NavigationPage: &gameGuidePage,
		settings:       settings,
		parent:         parent,

		toastOverlay: toastOverlay,
	}

	return &grp
}
