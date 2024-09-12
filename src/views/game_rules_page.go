package views

import (
	"github.com/tfuxu/floodit/src/constants"

	"github.com/diamondburned/gotk4-adwaita/pkg/adw"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

type GameRulesPage struct {
	*adw.Bin
	settings *gio.Settings
	parent   *MainWindow

	toastOverlay *adw.ToastOverlay
}

func NewGameRulesPage(parent *MainWindow, settings *gio.Settings, toastOverlay *adw.ToastOverlay) *GameRulesPage {
	builder := gtk.NewBuilderFromResource(constants.RootPath + "/ui/game_rules_page.ui")

	gameRulesPage := builder.GetObject("game_rules_page").Cast().(*adw.Bin)

	grp := GameRulesPage{
		Bin:      gameRulesPage,
		settings: settings,
		parent:   parent,

		toastOverlay: toastOverlay,
	}

	return &grp
}
