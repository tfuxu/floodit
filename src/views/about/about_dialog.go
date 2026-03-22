package about

import (
	"github.com/tfuxu/floodit/src/constants"

	"codeberg.org/puregotk/puregotk/v4/adw"
	"codeberg.org/puregotk/puregotk/v4/gtk"
	. "github.com/pojntfx/go-gettext/pkg/i18n"
)

type AboutDialog struct {
	*adw.AboutDialog
}

func NewAboutDialog() *AboutDialog {
	dialog := adw.NewAboutDialog()
	dialog.SetApplicationName(L("Flood It"))
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

	// TRANSLATORS: This is a place to put your credits (formats: "Name https://example.com" or "Name <email@example.com>", no quotes) and is not meant to be translated literally.
	dialog.SetTranslatorCredits(L("translator-credits"))
	dialog.SetCopyright(L("Copyright © 2026 tfuxu"))
	dialog.SetLicenseType(gtk.LicenseGpl30Value)
	dialog.SetVersion(constants.Version)
	dialog.SetReleaseNotesVersion(constants.RelVer)

	a := AboutDialog{
		AboutDialog: dialog,
	}

	return &a
}
