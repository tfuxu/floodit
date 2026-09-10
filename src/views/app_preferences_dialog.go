package views

import (
	"github.com/tfuxu/floodit/src/constants"
	
	"codeberg.org/puregotk/puregotk/v4/adw"
	"codeberg.org/puregotk/puregotk/v4/gio"
	"codeberg.org/puregotk/puregotk/v4/gtk"
)

type AppPreferencesDialog struct {
	adw.PreferencesDialog
	settings *gio.Settings
	parent   *MainWindow

	showNumbersSwitch *adw.SwitchRow
}

func NewAppPreferencesDialog(parent *MainWindow, settings *gio.Settings) *AppPreferencesDialog {
	builder := gtk.NewBuilderFromResource(constants.RootPath + "/ui/app_preferences_dialog.ui")

	var appPreferencesDialog adw.PreferencesDialog
	builder.GetObject("app_preferences_dialog").Cast(&appPreferencesDialog)
	defer appPreferencesDialog.Unref()

	var showNumbersSwitch adw.SwitchRow
	builder.GetObject("show_numbers_switch").Cast(&showNumbersSwitch)
	defer showNumbersSwitch.Unref()

	apd := AppPreferencesDialog{
		PreferencesDialog: appPreferencesDialog,
		settings:          settings,
		parent:            parent,

		showNumbersSwitch: &showNumbersSwitch,
	}

	apd.setupBindings()
	apd.showNumbersSwitch.SetActive(settings.GetBoolean("show-color-numbers"))

	return &apd
}

func (apd *AppPreferencesDialog) setupBindings() {
	apd.settings.Bind("show-color-numbers", &apd.showNumbersSwitch.Object, "active", gio.GSettingsBindDefaultValue)
}
