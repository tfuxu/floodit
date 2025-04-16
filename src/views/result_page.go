package views

import (
	"fmt"

	"github.com/tfuxu/floodit/src/constants"

	"github.com/diamondburned/gotk4-adwaita/pkg/adw"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

// TODO: Make this translatable
var ResultStates = map[string]string{
	"win": "You Won!", "lose": "You Lost!",
}

type ResultPage struct {
	*adw.Bin
	settings *gio.Settings
	parent   *MainWindow

	resultState string

	toastOverlay *adw.ToastOverlay

	titleLabel       *gtk.Label
	descriptionLabel *gtk.Label
}

func NewResultPage(parent *MainWindow, settings *gio.Settings, toastOverlay *adw.ToastOverlay) *ResultPage {
	builder := gtk.NewBuilderFromResource(constants.RootPath + "/ui/result_page.ui")

	resultPage := builder.GetObject("result_page").Cast().(*adw.Bin)

	titleLabel := builder.GetObject("title_label").Cast().(*gtk.Label)
	descriptionLabel := builder.GetObject("description_label").Cast().(*gtk.Label)

	rp := ResultPage{
		Bin:      resultPage,
		settings: settings,
		parent:   parent,

		toastOverlay: toastOverlay,

		titleLabel:       titleLabel,
		descriptionLabel: descriptionLabel,
	}

	return &rp
}

func (r *ResultPage) SetResultState(isWin bool) {
	if isWin {
		r.resultState = "win"
	} else {
		r.resultState = "lose"
	}

	r.refresh()
}

func (r *ResultPage) refresh() {
	var classes []string
	var description string

	if r.resultState == "win" {
		classes = []string{"title-1", "success"}
		description = fmt.Sprintf(
			"You flooded the <b>%s</b> board in <b>%d</b> moves!", // TODO: Make this translatable
			r.parent.gamePage.board.Name,
			r.parent.gamePage.board.Step,
		)
	} else {
		classes = []string{"title-1", "error"}
		description = fmt.Sprintf(
			"You failed to finish the <b>%s</b> board.\nBetter luck next time!", // TODO: Make this translatable
			r.parent.gamePage.board.Name,
		)
	}

	r.titleLabel.SetLabel(ResultStates[r.resultState])
	r.descriptionLabel.SetLabel(description)
	r.titleLabel.SetCSSClasses(classes)
}
