package board

import (
	"log/slog"
	"runtime"
	"strconv"
	"unsafe"

	"github.com/tfuxu/floodit/src/constants"

	"github.com/tfuxu/floodit/src/backend"

	"codeberg.org/puregotk/puregotk/v4/gdk"
	"codeberg.org/puregotk/puregotk/v4/gio"
	"codeberg.org/puregotk/puregotk/v4/glib"
	"codeberg.org/puregotk/puregotk/v4/gobject"
	"codeberg.org/puregotk/puregotk/v4/graphene"
	"codeberg.org/puregotk/puregotk/v4/gsk"
	"codeberg.org/puregotk/puregotk/v4/gtk"
	"codeberg.org/puregotk/puregotk/v4/pango"
)

var gTypeGameBoard gobject.Type

type GameBoard struct {
	gtk.Widget

	settings *gio.Settings

	board *backend.Board
}

func NewGameBoard(board *backend.Board, settings *gio.Settings, FirstPropertyNameVar string, varArgs ...interface{}) GameBoard {
	object := gobject.NewObject(gTypeGameBoard, FirstPropertyNameVar, varArgs...)

	var v GameBoard
	object.Cast(&v)

	gb := (*GameBoard)(unsafe.Pointer(object.GetData(constants.DataKeyGoInstance)))
	gb.settings = settings
	gb.board = board

	return v
}

func init() {
	var gbClassInit gobject.ClassInitFunc = func(type_class *gobject.TypeClass, class_data uintptr) {
		objectClass := (*gobject.ObjectClass)(unsafe.Pointer(type_class))

		objectClass.OverrideConstructed(func(o *gobject.Object) {
			parentObjClass := (*gobject.ObjectClass)(unsafe.Pointer(type_class.PeekParent()))
			parentObjClass.GetConstructed()(o)

			var parent gtk.Widget
			o.Cast(&parent)

			gb := &GameBoard{
				Widget: parent,
			}

			//gb.SetOverflow(gtk.OverflowHiddenValue)

			var pinner runtime.Pinner
			pinner.Pin(gb)

			var cleanupCallback glib.DestroyNotify = func(data uintptr) {
				pinner.Unpin()
			}
			o.SetDataFull(constants.DataKeyGoInstance, uintptr(unsafe.Pointer(gb)), &cleanupCallback)
		})

		widgetClass := (*gtk.WidgetClass)(unsafe.Pointer(type_class))

		widgetClass.OverrideSnapshot(func(widget *gtk.Widget, snapshot *gtk.Snapshot) {
			gb := (*GameBoard)(unsafe.Pointer(widget.GetData(constants.DataKeyGoInstance)))
			if gb == nil {
				return
			}

			slog.Debug("GameBoard snapshot hit!")

			if gb.board != nil {
				width := int(widget.GetWidth())
				height := int(widget.GetHeight())

				boardMatrix := gb.board.Matrix
				boardRows := gb.board.Rows
				boardCols := gb.board.Columns

				rectWidth := width / boardCols
				rectHeight := height / boardRows
				xOffset := (width - rectWidth*boardCols) / 2
				yOffset := (height - rectHeight*boardRows) / 2

				pangoContext := widget.CreatePangoContext()

				snapshot.Save()

				roundedRect := gsk.RoundedRect{}
				roundedRect.InitFromRect(
					graphene.RectAlloc().Init(
						float32(xOffset),
						float32(yOffset),
						float32(rectWidth*boardCols),
						float32(rectHeight*boardRows),
					),
					12.0,
				)

				snapshot.PushRoundedClip(&roundedRect)

				for row := 0; row < boardRows; row++ {
					for col := 0; col < boardCols; col++ {
						x := rectWidth*col + xOffset
						y := rectHeight*row + yOffset
						var hexCode string
						var colorLabel string

						// TODO: This is a dirty workaround to get this working with array.
						// Make sure in future to use array indexes in game matrix instead.
						for i, color := range backend.DefaultColors {
							if color[0] == boardMatrix[row][col] {
								hexCode = color[1]
								colorLabel = strconv.Itoa(i + 1)
							}
						}

						color := gdk.RGBA{}
						if ok := color.Parse(hexCode); !ok {
							// TODO: Show user some feedback in UI when this happens
							slog.Error("Failed to convert hex values to Cairo compatible RGB channels.")
							return
						}

						snapshot.AppendColor(
							&color,
							graphene.RectAlloc().Init(
								float32(x),
								float32(y),
								float32(rectWidth),
								float32(rectHeight),
							),
						)

						if gb.settings.GetBoolean("show-color-numbers") {
							// TODO: Check how to get what font is currently used for UI
							fontDescription := pango.FontDescriptionFromString(
								"Adwaita Sans Bold " + strconv.Itoa(rectWidth/2),
							)

							var layoutWidth int32
							var layoutHeight int32

							layout := pango.NewLayout(pangoContext)
							layout.SetFontDescription(fontDescription)
							layout.SetText(colorLabel, -1)
							layout.GetPixelSize(&layoutWidth, &layoutHeight)

							// TODO: Switch between black and white depending on a background/text contrast
							black := gdk.RGBA{
								Red:   0.0,
								Green: 0.0,
								Blue:  0.0,
								Alpha: 1.0,
							}

							centerX := x + (rectWidth-int(layoutWidth))/2
							centerY := y + (rectHeight-int(layoutHeight))/2

							snapshot.Save()

							snapshot.Translate(
								graphene.PointAlloc().Init(
									float32(centerX),
									float32(centerY),
								),
							)
							snapshot.AppendLayout(layout, &black)

							snapshot.Restore()
						}
					}
				}

				snapshot.Pop()
				snapshot.Restore()
			}
		})
	}

	var gbInstanceInit gobject.InstanceInitFunc = func(type_instance *gobject.TypeInstance, type_class *gobject.TypeClass) {}

	var gbParentQuery gobject.TypeQuery
	gobject.NewTypeQuery(gtk.WidgetGLibType(), &gbParentQuery)

	gTypeGameBoard = gobject.TypeRegisterStaticSimple(
		gbParentQuery.Type,
		"GameBoard",
		gbParentQuery.ClassSize,
		&gbClassInit,
		gbParentQuery.InstanceSize,
		&gbInstanceInit,
		0,
	)
}
