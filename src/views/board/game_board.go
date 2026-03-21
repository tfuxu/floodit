package board

import (
	"log/slog"
	"runtime"
	"unsafe"

	"github.com/tfuxu/floodit/src/constants"

	"github.com/tfuxu/floodit/src/backend"
	"github.com/tfuxu/floodit/src/backend/utils"

	"codeberg.org/puregotk/puregotk/v4/gdk"
	"codeberg.org/puregotk/puregotk/v4/glib"
	"codeberg.org/puregotk/puregotk/v4/gobject"
	"codeberg.org/puregotk/puregotk/v4/graphene"
	"codeberg.org/puregotk/puregotk/v4/gtk"
)

var gTypeGameBoard gobject.Type

type GameBoard struct {
	gtk.Widget

	board *backend.Board
}

func NewGameBoard(board *backend.Board, FirstPropertyNameVar string, varArgs ...interface{}) GameBoard {
	object := gobject.NewObject(gTypeGameBoard, FirstPropertyNameVar, varArgs...)

	var v GameBoard
	object.Cast(&v)

	gb := (*GameBoard)(unsafe.Pointer(object.GetData(constants.DataKeyGoInstance)))
	gb.board = board

	return v
}

// TODO: Implement board rounding
func (gb *GameBoard) roundedRect(x, y, width, height, cornerRadius float64) {
	/*ctx.NewSubPath()
	ctx.Arc(x+width-cornerRadius, y+cornerRadius, cornerRadius, -math.Pi/2, 0)
	ctx.Arc(x+width-cornerRadius, y+height-cornerRadius, cornerRadius, 0, math.Pi/2)
	ctx.Arc(x+cornerRadius, y+height-cornerRadius, cornerRadius, math.Pi/2, math.Pi)
	ctx.Arc(x+cornerRadius, y+cornerRadius, cornerRadius, math.Pi, 3*math.Pi/2)
	ctx.ClosePath()*/
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
			slog.Debug("Current board addr:", gb.board)

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

				snapshot.Save()

				for row := 0; row < boardRows; row++ {
					for col := 0; col < boardCols; col++ {
						x := rectWidth*col + xOffset
						y := rectHeight*row + yOffset
						var hexCode string

						// TODO: This is a dirty workaround to get this working with array.
						// Make sure in future to use array indexes in game matrix instead.
						for _, color := range backend.DefaultColors {
							if color[0] == boardMatrix[row][col] {
								hexCode = color[1]
							}
						}

						cairoRGB, err := utils.HexToCairoRGB(hexCode)
						if err != nil {
							// TODO: Show user some feedback in UI when this happens
							slog.Error("Failed to convert hex values to Cairo compatible RGB channels:", "msg", err)
						}

						color := gdk.RGBA{
							Red:   float32(cairoRGB[0]),
							Green: float32(cairoRGB[1]),
							Blue:  float32(cairoRGB[2]),
							Alpha: 1.0,
						}

						snapshot.AppendColor(
							&color,
							graphene.RectAlloc().Init(float32(x), float32(y), float32(rectWidth), float32(rectHeight)),
						)
					}
				}

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
