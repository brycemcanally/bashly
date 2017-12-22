/*
Package views implements the functionality for setting up and managing
simple views with functionality that can be used inside of boxes.
*/
package views

import "github.com/jroimartin/gocui"

// View is a view that can be used within a box.
type View interface {
	Set(gui *gocui.Gui, x0, y0, x1, y1 int) error
	Delete(gui *gocui.Gui) error
}
