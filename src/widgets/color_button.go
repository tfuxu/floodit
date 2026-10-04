package widgets

import (
	"runtime"
	"unsafe"

	"github.com/tfuxu/floodit/src/constants"

	"codeberg.org/puregotk/puregotk/v4/gdk"
	"codeberg.org/puregotk/puregotk/v4/glib"
	"codeberg.org/puregotk/puregotk/v4/gobject"
	"codeberg.org/puregotk/puregotk/v4/graphene"
	"codeberg.org/puregotk/puregotk/v4/gtk"
)

var gTypeColorButton gobject.Type

type ColorButton struct {
	gtk.Button

	color *gdk.RGBA
}

func NewColorButton(color *gdk.RGBA) *ColorButton {
	object := gobject.NewObject(gTypeColorButton, "css-name")

	var v ColorButton
	object.Cast(&v)

	cb := (*ColorButton)(unsafe.Pointer(object.GetData(constants.DataKeyGoInstance)))
	cb.color = color

	return &v
}

func init() {
	var cbClassInit gobject.ClassInitFunc = func(type_class *gobject.TypeClass, class_data uintptr) {
		objectClass := (*gobject.ObjectClass)(unsafe.Pointer(type_class))

		objectClass.OverrideConstructed(func(object *gobject.Object) {
			parentObjClass := (*gobject.ObjectClass)(unsafe.Pointer(type_class.PeekParent()))
			parentObjClass.GetConstructed()(object)

			var parent gtk.Button
			object.Cast(&parent)

			cb := &ColorButton{
				Button: parent,
			}

			// Set button properties
			cb.SetOverflow(gtk.OverflowHiddenValue)
			cb.SetCssClasses([]string{"card", "circular", "color-button"})

			var pinner runtime.Pinner
			pinner.Pin(cb)

			var cleanupCallback glib.DestroyNotify = func(data uintptr) {
				pinner.Unpin()
			}
			object.SetDataFull(constants.DataKeyGoInstance, uintptr(unsafe.Pointer(cb)), &cleanupCallback)
		})

		widgetClass := (*gtk.WidgetClass)(unsafe.Pointer(type_class))

		widgetClass.OverrideSnapshot(func(widget *gtk.Widget, snapshot *gtk.Snapshot) {
			cb := (*ColorButton)(unsafe.Pointer(widget.GetData(constants.DataKeyGoInstance)))
			if cb == nil {
				return
			}

			widgetWidth := widget.GetWidth()
			widgetHeight := widget.GetHeight()
			color := cb.color

			bounds := graphene.RectAlloc().Init(
				0.0,
				0.0,
				float32(widgetWidth),
				float32(widgetHeight),
			)
			defer bounds.Free()

			snapshot.AppendColor(color, bounds)

			// Check if button has a child (pointer will be zero if no child is set)
			if cb.GetChild().GoPointer() != 0 {
				// Snapshot child to make it visible above the node
				cb.SnapshotChild(cb.GetChild(), snapshot)
			}
		})
	}

	var cbInstanceInit gobject.InstanceInitFunc = func(type_instance *gobject.TypeInstance, type_class *gobject.TypeClass) {}

	var cbParentQuery gobject.TypeQuery
	gobject.NewTypeQuery(gtk.ButtonGLibType(), &cbParentQuery)

	gTypeColorButton = gobject.TypeRegisterStaticSimple(
		cbParentQuery.Type,
		"ColorButton",
		cbParentQuery.ClassSize,
		&cbClassInit,
		cbParentQuery.InstanceSize,
		&cbInstanceInit,
		0,
	)
}
