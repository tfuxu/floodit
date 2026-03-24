package keyboard

import (
	"fmt"
	"strings"

	"github.com/tfuxu/floodit/src/backend/utils"
	"github.com/tfuxu/floodit/src/constants"

	"codeberg.org/puregotk/puregotk/v4/gdk"
	"codeberg.org/puregotk/puregotk/v4/gtk"
)

type ColorKeyboard struct {
	*gtk.Box

	cssProvider *gtk.CssProvider

	rowFirst  *gtk.Box
	rowSecond *gtk.Box
}

// NewColorKeyboard creates a new instance of ColorKeyboard.
// It takes currently used color palette.
func NewColorKeyboard(colorPalette [][2]string) *ColorKeyboard {
	builder := gtk.NewBuilderFromResource(constants.RootPath + "/ui/color_keyboard.ui")

	var keyboard gtk.Box
	builder.GetObject("color_keyboard").Cast(&keyboard)
	defer keyboard.Unref()

	var rowFirst gtk.Box
	builder.GetObject("row_first").Cast(&rowFirst)
	defer rowFirst.Unref()

	var rowSecond gtk.Box
	builder.GetObject("row_second").Cast(&rowSecond)
	defer rowSecond.Unref()

	cssProvider := gtk.NewCssProvider()

	gtk.StyleContextAddProviderForDisplay(
		gdk.DisplayGetDefault(),
		cssProvider,
		uint32(gtk.STYLE_PROVIDER_PRIORITY_USER+1),
	)

	ck := ColorKeyboard{
		Box: &keyboard,

		cssProvider: cssProvider,

		rowFirst:  &rowFirst,
		rowSecond: &rowSecond,
	}

	ck.setupButtons(colorPalette)

	return &ck
}

func (ck *ColorKeyboard) setupButtons(colorPalette [][2]string) {
	buttonStore := make([]*gtk.Button, len(colorPalette))

	var buttonColors []string

	for i, color := range colorPalette {
		colorName := color[0]
		colorHex := color[1]

		button := gtk.NewButton()
		button.SetTooltipText(utils.ToSentenceString(colorName))
		buttonColors = append(buttonColors, fmt.Sprintf(".%s-button { background-color: %s; }", colorName, colorHex))
		button.SetCssClasses([]string{"card", "circular", "color-button", fmt.Sprintf("%s-button", colorName)})
		button.SetActionName("game.select-color")
		button.SetActionTarget("s", colorName)

		buttonStore[i] = button
	}

	ckCSS := strings.Join(buttonColors, " ")
	ck.cssProvider.LoadFromString(ckCSS)

	colorNo := 1
	currentRow := ck.rowFirst
	for _, button := range buttonStore {
		if colorNo > 4 {
			currentRow = ck.rowSecond
			colorNo = 1
		}

		currentRow.Append(&button.Widget)
		colorNo += 1
	}
}
