package views

import (
	"runtime"
	"unsafe"

	"codeberg.org/puregotk/puregotk/v4/adw"
	"codeberg.org/puregotk/puregotk/v4/glib"
	"codeberg.org/puregotk/puregotk/v4/gobject"
	"codeberg.org/puregotk/puregotk/v4/gtk"
	"github.com/tfuxu/floodit/src/constants"
)

var gTypeGameRulesDialog gobject.Type

type GameRulesDialog struct {
	adw.Dialog
}

func NewGameRulesDialog(FirstPropertyNameVar string, varArgs ...interface{}) *GameRulesDialog {
	object := gobject.NewObject(gTypeGameRulesDialog, FirstPropertyNameVar, varArgs...)

	var v GameRulesDialog
	object.Cast(&v)

	return &v
}

func init() {
	var grdClassInit gobject.ClassInitFunc = func(type_class *gobject.TypeClass, class_data uintptr) {
		typeClass := (*gtk.WidgetClass)(unsafe.Pointer(type_class))
		typeClass.SetTemplateFromResource(constants.RootPath + "/ui/game_rules_dialog.ui")

		objectClass := (*gobject.ObjectClass)(unsafe.Pointer(type_class))

		objectClass.OverrideConstructed(func(o *gobject.Object) {
			parentObjClass := (*gobject.ObjectClass)(unsafe.Pointer(type_class.PeekParent()))
			parentObjClass.GetConstructed()(o)

			var parent adw.Dialog
			o.Cast(&parent)

			parent.InitTemplate()

			grd := &GameRulesDialog{
				Dialog: parent,
			}

			var pinner runtime.Pinner
			pinner.Pin(grd)

			var cleanupCallback glib.DestroyNotify = func(data uintptr) {
				pinner.Unpin()
			}
			o.SetDataFull(constants.DataKeyGoInstance, uintptr(unsafe.Pointer(grd)), &cleanupCallback)
		})
	}

	var grdInstanceInit gobject.InstanceInitFunc = func(type_instance *gobject.TypeInstance, type_class *gobject.TypeClass) {}

	var grdParentQuery gobject.TypeQuery
	gobject.NewTypeQuery(adw.DialogGLibType(), &grdParentQuery)

	gTypeGameRulesDialog = gobject.TypeRegisterStaticSimple(
		grdParentQuery.Type,
		"GameRulesDialog",
		grdParentQuery.ClassSize,
		&grdClassInit,
		grdParentQuery.InstanceSize,
		&grdInstanceInit,
		0,
	)
}
