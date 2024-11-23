package keyboard

import (
	"fmt"
	"strings"

	"github.com/tfuxu/floodit/src/constants"

	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

type ColorKeyboard struct {
	*gtk.Box

	cssProvider *gtk.CSSProvider

	rowFirst  *gtk.Box
	rowSecond *gtk.Box
}

// NewColorKeyboard creates a new instance of ColorKeyboard.
// It takes currently used color palette and a function to call when the button is pressed.
func NewColorKeyboard(colorPalette [][2]string, callback func(colorName string)) *ColorKeyboard {
	builder := gtk.NewBuilderFromResource(constants.RootPath + "/ui/color_keyboard.ui")

	keyboard := builder.GetObject("color_keyboard").Cast().(*gtk.Box)

	rowFirst := builder.GetObject("row_first").Cast().(*gtk.Box)
	rowSecond := builder.GetObject("row_second").Cast().(*gtk.Box)

	cssProvider := gtk.NewCSSProvider()

	gtk.StyleContextAddProviderForDisplay(
		gdk.DisplayGetDefault(),
		cssProvider,
		gtk.STYLE_PROVIDER_PRIORITY_USER+1,
	)

	ck := ColorKeyboard{
		Box: keyboard,

		cssProvider: cssProvider,

		rowFirst:  rowFirst,
		rowSecond: rowSecond,
	}

	ck.setupButtons(colorPalette, callback)

	return &ck
}

func (ck *ColorKeyboard) setupButtons(colorPalette [][2]string, callback func(colorName string)) {
	buttonStore := make([]*gtk.Button, len(colorPalette))

	var buttonColors []string

	for i, color := range colorPalette {
		colorName := color[0]
		colorHex := color[1]

		button := gtk.NewButton()
		button.SetTooltipText(colorName) // TODO: Capitalize text
		buttonColors = append(buttonColors, fmt.Sprintf(".%s-button { background-color: %s; }", colorName, colorHex))
		button.SetCSSClasses([]string{"card", "circular", "color-button", fmt.Sprintf("%s-button", colorName)})

		button.ConnectClicked(func() {
			callback(colorName)
		})

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

		currentRow.Append(button)
		colorNo += 1
	}
}
