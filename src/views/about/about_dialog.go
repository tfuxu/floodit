package about

import (
	"github.com/tfuxu/floodit/src/constants"

	"codeberg.org/puregotk/puregotk/v4/adw"
	"codeberg.org/puregotk/puregotk/v4/gtk"
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
	dialog.SetSupportUrl(constants.HelpUrl)
	dialog.SetIssueUrl(constants.BugtrackerUrl)

	dialog.SetDevelopers([]string{
		"tfuxu https://github.com/tfuxu",
	})
	dialog.SetDesigners([]string{
		"tfuxu https://github.com/tfuxu",
	})

	dialog.SetTranslatorCredits("translator-credits")
	dialog.SetCopyright("Copyright © 2024-2025 tfuxu")
	dialog.SetLicenseType(gtk.LicenseGpl30Value)
	dialog.SetVersion(constants.Version)
	dialog.SetReleaseNotesVersion(constants.RelVer)

	a := AboutDialog{
		AboutDialog: dialog,
	}

	return &a
}
