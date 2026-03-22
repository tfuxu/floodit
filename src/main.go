package main

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/tfuxu/floodit/src/constants"
	"github.com/tfuxu/floodit/src/views"
	"github.com/tfuxu/floodit/src/views/about"

	"codeberg.org/puregotk/puregotk/v4/adw"
	"codeberg.org/puregotk/puregotk/v4/gio"

	"github.com/pojntfx/go-gettext/pkg/i18n"
)

func init() {
	// Select Logger level depending on current build type
	if constants.BuildType == "debug" {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	} else {
		slog.SetLogLoggerLevel(slog.LevelInfo)
	}

	// TODO: Check if this will be needed after devenv removal
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

	// Initialize i18n
	if err := i18n.InitI18n("floodit", constants.LocaleDir, slog.Default()); err != nil {
		slog.Error("Failed to initialize i18n.", "msg", err)
		os.Exit(1)
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

	app := adw.NewApplication(constants.AppID, gio.GApplicationDefaultFlagsValue)
	defer app.Unref()

	app.SetResourceBasePath(constants.RootPath)
	setupActions(app)

	app.ConnectActivate(new(func(gio.Application) {
		doActivate(app, settings)
	}))

	if code := app.Run(int32(len(os.Args)), os.Args); code > 0 {
		os.Exit(int(code))
	}
}

func doActivate(app *adw.Application, settings *gio.Settings) {
	window := views.NewMainWindow(app, settings)
	window.Present()
}

func setupActions(app *adw.Application) {
	aboutAction := gio.NewSimpleAction("about", nil)
	aboutAction.ConnectActivate(new(func(gio.SimpleAction, uintptr) {
		onAbout(app)
	}))
	app.AddAction(aboutAction)

	quitAction := gio.NewSimpleAction("quit", nil)
	quitAction.ConnectActivate(new(func(gio.SimpleAction, uintptr) {
		onQuit(app)
	}))
	app.AddAction(quitAction)

	app.SetAccelsForAction("win.play-again", []string{"<Primary>R"})
	app.SetAccelsForAction("win.show-game-select", []string{"<Primary>N"})
	app.SetAccelsForAction("app.game-rules", []string{"F1"})
	app.SetAccelsForAction("app.quit", []string{"<Primary>Q"})
}

func onAbout(app *adw.Application) {
	a := about.NewAboutDialog()
	a.Present(&app.GetActiveWindow().Widget)
}

func onQuit(app *adw.Application) {
	app.Quit()
}
