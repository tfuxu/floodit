package keyboard

import (
	"strconv"
	"log/slog"

	"github.com/tfuxu/floodit/src/backend/utils"
	"github.com/tfuxu/floodit/src/constants"
	"github.com/tfuxu/floodit/src/widgets"

	"codeberg.org/puregotk/puregotk/v4/gdk"
	"codeberg.org/puregotk/puregotk/v4/gio"
	"codeberg.org/puregotk/puregotk/v4/gtk"
	"codeberg.org/puregotk/puregotk/v4/pango"
)

type ColorKeyboard struct {
	*gtk.Box
	settings *gio.Settings

	buttonStore []*widgets.ColorButton

	rowFirst  *gtk.Box
	rowSecond *gtk.Box
}

// NewColorKeyboard creates a new instance of ColorKeyboard.
// It takes a pointer to the gio.Settings object and a currently used color palette.
func NewColorKeyboard(settings *gio.Settings, colorPalette [][2]string) *ColorKeyboard {
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

	ck := ColorKeyboard{
		Box:         &keyboard,
		settings:    settings,

		buttonStore: make([]*widgets.ColorButton, len(colorPalette)),

		rowFirst:  &rowFirst,
		rowSecond: &rowSecond,
	}

	ck.setupButtons(colorPalette)
	ck.setupSignals()

	return &ck
}

func (ck *ColorKeyboard) setupSignals() {
	ck.settings.ConnectChanged(new(func(settings gio.Settings, key string) {
		if key == "show-color-numbers" {
			ck.setColorNumbers(settings.GetBoolean("show-color-numbers"))
		}
	}))
}

func (ck *ColorKeyboard) setupButtons(colorPalette [][2]string) {
	for i, color := range colorPalette {
		colorName := color[0]
		colorHex := color[1]
		buttonLabel := strconv.Itoa(i + 1)

		label := gtk.NewLabel(buttonLabel)
		label.SetVisible(ck.settings.GetBoolean("show-color-numbers"))
		label.AddCssClass("title-1")
		label.SetHalign(gtk.AlignCenterValue)
		label.SetValign(gtk.AlignCenterValue)

		// TODO: Add text color value to DefaultColors to remove this jank
		if buttonLabel == "1" || buttonLabel == "5" || buttonLabel == "6" {
			label.SetAttributes(pango.AttrListFromString("foreground white"))
		} else {
			label.SetAttributes(pango.AttrListFromString("foreground black"))
		}

		color := gdk.RGBA{}
		if ok := color.Parse(colorHex); !ok {
			// TODO: Show user some feedback in UI when this happens
			slog.Error("Failed to convert hex values to Cairo compatible RGB channels.", "colorName", colorName, "colorHex", colorHex)
		}

		button := widgets.NewColorButton(&color)
		button.SetChild(&label.Widget)
		button.SetTooltipText(utils.ToSentenceString(colorName))
		button.SetActionName("game.select-color")
		button.SetActionTarget("s", colorName)

		ck.buttonStore[i] = button
	}

	colorNo := 1
	currentRow := ck.rowFirst
	for _, button := range ck.buttonStore {
		if colorNo > 4 {
			currentRow = ck.rowSecond
			colorNo = 1
		}

		currentRow.Append(&button.Widget)
		colorNo += 1
	}
}

func (ck *ColorKeyboard) setColorNumbers(showColorNumbers bool) {
	for _, button := range ck.buttonStore {
		button.GetChild().SetVisible(showColorNumbers)
	}
}
