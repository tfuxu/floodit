package views

import (
	"github.com/tfuxu/floodit/src/constants"

	"codeberg.org/puregotk/puregotk/v4/adw"
	"codeberg.org/puregotk/puregotk/v4/gio"
	"codeberg.org/puregotk/puregotk/v4/gtk"
)

type PreferencesDialog struct {
	adw.PreferencesDialog
	settings *gio.Settings
	parent   *MainWindow

	showNumbersSwitch *adw.SwitchRow
}

func NewPreferencesDialog(parent *MainWindow, settings *gio.Settings) *PreferencesDialog {
	builder := gtk.NewBuilderFromResource(constants.RootPath + "/ui/preferences_dialog.ui")

	var preferencesDialog adw.PreferencesDialog
	builder.GetObject("preferences_dialog").Cast(&preferencesDialog)
	defer preferencesDialog.Unref()

	var showNumbersSwitch adw.SwitchRow
	builder.GetObject("show_numbers_switch").Cast(&showNumbersSwitch)
	defer showNumbersSwitch.Unref()

	pd := PreferencesDialog{
		PreferencesDialog: preferencesDialog,
		settings:          settings,
		parent:            parent,

		showNumbersSwitch: &showNumbersSwitch,
	}

	pd.setupBindings()
	pd.showNumbersSwitch.SetActive(settings.GetBoolean("show-color-numbers"))

	return &pd
}

func (pd *PreferencesDialog) setupBindings() {
	pd.settings.Bind("show-color-numbers", &pd.showNumbersSwitch.Object, "active", gio.GSettingsBindDefaultValue)
}
