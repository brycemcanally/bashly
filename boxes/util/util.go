package util

import (
	"github.com/jroimartin/gocui"
)

// ScrollDown scrolls the content of a view down.
func ScrollDown(_ *gocui.Gui, view *gocui.View) error {
	x, y := view.Origin()
	_, maxY := view.Size()

	if _, err := view.Line(maxY + 1); err == nil {
		y++
		view.SetOrigin(x, y)
	}

	return nil
}

// ScrollUp scrolls the content of a view up.
func ScrollUp(_ *gocui.Gui, view *gocui.View) error {
	x, y := view.Origin()
	if y > 0 {
		y--
		view.SetOrigin(x, y)
	}

	return nil
}

// PositionIndex gets the index in the string representation of the view's buffer of a given position.
func PositionIndex(view *gocui.View, x, y int) int {
	_, oy := view.Origin()
	y += oy
	idx := 0

	width, _ := view.Size()
	tmpX := 0

	for _, val := range view.Buffer() {
		if x <= 0 && y <= 0 {
			break
		}

		if tmpX >= width {
			y--
			tmpX = 0
			continue
		}
		tmpX++

		if val == '\n' {
			y--
			tmpX = 0
		} else if y == 0 {
			x--
		}
		idx++
	}

	return idx
}

// IndexPosition gets the position of a given index in the string representation of the view's buffer.
func IndexPosition(view *gocui.View, index int) (int, int) {
	x, y := 0, 0

	width, _ := view.Size()

	for i, val := range view.Buffer() {
		if i >= index {
			break
		}
		if val == '\n' || x >= width {
			y++
			x = 0
		} else {
			x++
		}
	}

	return x, y
}

// Coordinates are a set of coordinates that represent the dimensions of a GUI element.
type Coordinates struct {
	X0, Y0, X1, Y1 int
}

// RealCoordinates computes the real coordinates on a GUI using a set of unit coordinates.
func RealCoordinates(gui *gocui.Gui, unit *Coordinates) (x0, y0, x1, y1 int) {
	maxX, maxY := gui.Size()

	x0 = maxX * unit.X0 / 100
	y0 = maxY * unit.Y0 / 100
	x1 = maxX*unit.X1/100 - 1
	y1 = maxY*unit.Y1/100 - 1

	return
}
