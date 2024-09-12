package views

import (
	"github.com/tfuxu/floodit/src/constants"

	"github.com/diamondburned/gotk4-adwaita/pkg/adw"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

type StartingView struct {
	*adw.NavigationView
	settings *gio.Settings
	parent   *MainWindow

	toastOverlay *adw.ToastOverlay

	difficultyPage *DifficultyPage
	customModePage *CustomModePage
}

func NewStartingView(parent *MainWindow, settings *gio.Settings, toastOverlay *adw.ToastOverlay) *StartingView {
	builder := gtk.NewBuilderFromResource(constants.RootPath + "/ui/starting_view.ui")

	startingView := builder.GetObject("starting_view").Cast().(*adw.NavigationView)

	sv := StartingView{
		NavigationView: startingView,
		settings:       settings,
		parent:         parent,

		toastOverlay: toastOverlay,
	}
	sv.difficultyPage = NewDifficultyPage(&sv, settings, toastOverlay)
	sv.customModePage = NewCustomModePage(&sv, settings, toastOverlay)

	sv.setupNavigation()

	return &sv
}

func (sv *StartingView) setupNavigation() {
	sv.Add(sv.difficultyPage.NavigationPage)
	sv.Add(sv.customModePage.NavigationPage)
}
