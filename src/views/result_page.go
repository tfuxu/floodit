package views

import (
	"fmt"

	"github.com/tfuxu/floodit/src/constants"

	"codeberg.org/puregotk/puregotk/v4/adw"
	"codeberg.org/puregotk/puregotk/v4/gio"
	"codeberg.org/puregotk/puregotk/v4/gtk"
	. "github.com/pojntfx/go-gettext/pkg/i18n"
)

var ResultStates = map[string]string{
	"win": L("You Won!"), "lose": L("You Lost!"),
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

	var resultPage adw.Bin
	builder.GetObject("result_page").Cast(&resultPage)
	defer resultPage.Unref()

	var titleLabel gtk.Label
	builder.GetObject("title_label").Cast(&titleLabel)
	defer titleLabel.Unref()

	var descriptionLabel gtk.Label
	builder.GetObject("description_label").Cast(&descriptionLabel)
	defer descriptionLabel.Unref()

	rp := ResultPage{
		Bin:      &resultPage,
		settings: settings,
		parent:   parent,

		toastOverlay: toastOverlay,

		titleLabel:       &titleLabel,
		descriptionLabel: &descriptionLabel,
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
			// TRANSLATORS: DO NOT translate '<b>%s</b>' and '<b>%d</b>'.
			L("You flooded the <b>%s</b> board in <b>%d</b> moves!"),
			r.parent.gamePage.board.Name,
			r.parent.gamePage.board.Step,
		)
	} else {
		classes = []string{"title-1", "error"}
		description = fmt.Sprintf(
			// TRANSLATORS: DO NOT translate '<b>%s</b>' and '\n'.
			L("You failed to finish the <b>%s</b> board.\nBetter luck next time!"),
			r.parent.gamePage.board.Name,
		)
	}

	r.titleLabel.SetLabel(ResultStates[r.resultState])
	r.descriptionLabel.SetLabel(description)
	r.titleLabel.SetCssClasses(classes)
}
