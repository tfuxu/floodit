package about

import (
	"github.com/tfuxu/floodit/src/constants"

	"github.com/diamondburned/gotk4-adwaita/pkg/adw"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

type AboutDialog struct {
	*adw.AboutDialog
}

func NewAboutDialog() *AboutDialog {
	dialog := adw.NewAboutDialog()
	dialog.SetApplicationName("Flood It")
	dialog.SetApplicationIcon(constants.AppID)
	dialog.SetDeveloperName("tfuxu")
	dialog.SetWebsite(constants.ProjectUrl)
	dialog.SetSupportURL(constants.HelpUrl)
	dialog.SetIssueURL(constants.BugtrackerUrl)

	dialog.SetDevelopers([]string{
		"tfuxu https://github.com/tfuxu",
	})
	dialog.SetDesigners([]string{
		"tfuxu https://github.com/tfuxu",
	})

	dialog.SetTranslatorCredits("translator-credits")
	dialog.SetCopyright("Copyright © 2024 tfuxu")
	dialog.SetLicenseType(gtk.LicenseGPL30)
	dialog.SetVersion(constants.Version)
	dialog.SetReleaseNotesVersion(constants.RelVer)

	a := AboutDialog{
		AboutDialog: dialog,
	}

	return &a
}
