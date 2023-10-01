package views

import (
	//"github.com/tfuxu/flood_it/src/constants"

	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

type MainWindow struct {
	window *gtk.ApplicationWindow
	settings * gio.Settings
}

func NewMainWindow(app *gtk.Application, settings *gio.Settings) MainWindow {
	builder := gtk.NewBuilderFromResource("/io/github/tfuxu/flood_it/ui/main_window.ui")

	window := builder.GetObject("main_window").Cast().(*gtk.ApplicationWindow)
	window.SetApplication(app)
	window.ConnectUnrealize(func() {
		saveWindowProps(window, settings)
	})

	button := builder.GetObject("example_button").Cast().(*gtk.Button)
	button.ConnectClicked(func() {
		println("Hi ^_^")
	})

	return MainWindow{
		window: window,
		settings: settings,
	}
}

func saveWindowProps(window *gtk.ApplicationWindow, settings *gio.Settings) {
	width, height := window.DefaultSize()

	settings.SetInt("window-width", width)
	settings.SetInt("window-height", height)
	settings.SetBoolean("window-maximized", window.IsMaximized())
}
