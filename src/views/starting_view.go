package views

import (
	"github.com/tfuxu/floodit/src/constants"

	"codeberg.org/puregotk/puregotk/v4/adw"
	"codeberg.org/puregotk/puregotk/v4/gio"
	"codeberg.org/puregotk/puregotk/v4/gtk"
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

	var startingView adw.NavigationView
	builder.GetObject("starting_view").Cast(&startingView)
	defer startingView.Unref()

	sv := StartingView{
		NavigationView: &startingView,
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
