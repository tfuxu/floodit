package main

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/tfuxu/floodit/src/constants"
	"github.com/tfuxu/floodit/src/views"
	"github.com/tfuxu/floodit/src/views/about"

	"github.com/diamondburned/gotk4-adwaita/pkg/adw"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/glib/v2"
)

func init() {
	// Select Logger level depending on current build type
	if constants.BuildType == "debug" {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	} else {
		slog.SetLogLoggerLevel(slog.LevelInfo)
	}

	// Set `XDG_DATA_DIRS` to <builddir>/data if running in devenv
	if os.Getenv("MESON_DEVENV") == "1" {
		var data_dirs string

		if len(os.Getenv("XDG_DATA_DIRS")) != 0 {
			data_dirs = os.Getenv("XDG_DATA_DIRS")
		} else {
			data_dirs = "/usr/local/share/:/usr/share/"
		}

		os.Setenv("XDG_DATA_DIRS", fmt.Sprintf("%s:%s", constants.DataDir, data_dirs))
	}

	// Load resources
	resources, error := gio.ResourceLoad(filepath.Join(constants.PkgDataDir, "floodit.gresource"))
	if error != nil {
		slog.Error(error.Error())
		os.Exit(1)
	}
	gio.ResourcesRegister(resources)
}

func main() {
	settings := gio.NewSettings(constants.AppID)

	app := adw.NewApplication(constants.AppID, gio.ApplicationFlagsNone)
	app.SetResourceBasePath(constants.RootPath)

	setupActions(app)

	app.ConnectActivate(func() {
		doActivate(app, settings)
	})

	if code := app.Run(os.Args); code > 0 {
		os.Exit(code)
	}
}

func doActivate(app *adw.Application, settings *gio.Settings) {
	window := views.NewMainWindow(app, settings)
	window.Present()
}

func setupActions(app *adw.Application) {
	aboutAction := gio.NewSimpleAction("about", nil)
	aboutAction.ConnectActivate(func(parameter *glib.Variant) {
		onAbout(app)
	})
	app.AddAction(aboutAction)

	quitAction := gio.NewSimpleAction("quit", nil)
	quitAction.ConnectActivate(func(parameter *glib.Variant) {
		onQuit(app)
	})
	app.AddAction(quitAction)

	app.SetAccelsForAction("win.show-help-overlay", []string{"<Primary>question"})
	app.SetAccelsForAction("win.play-again", []string{"<Primary>R"})
	app.SetAccelsForAction("win.show-game-select", []string{"<Primary>N"})
	//app.SetAccelsForAction("win.show-game-rules", []string{"F1"})
	app.SetAccelsForAction("app.quit", []string{"<Primary>Q"})
}

func onAbout(app *adw.Application) {
	a := about.NewAboutDialog()
	a.Present(app.ActiveWindow())
}

func onQuit(app *adw.Application) {
	app.Quit()
}
