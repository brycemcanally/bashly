package boxes

import (
	"github.com/bryce/bashly/boxes/util"
	"github.com/bryce/bashly/cmds"
	"github.com/jroimartin/gocui"
)

// Script is a type of box that holds the editable script text.
type Script struct {
	name    string
	unit    util.Coordinates
	tabSize int
	command string
}

// NewScript creates a new script box.
func NewScript(cfg *Config) *Script {
	box := &Script{}
	box.name = cfg.Name
	box.unit.X0 = cfg.X0
	box.unit.Y0 = cfg.Y0
	box.unit.X1 = cfg.X1
	box.unit.Y1 = cfg.Y1
	box.tabSize = cfg.TabSize

	return box
}

// Name returns the name associated with this box.
func (box *Script) Name() string {
	return box.name
}

// Setup sets the keybindings for this box.
func (box *Script) Setup(gui *gocui.Gui, boxs *Boxes) error {
	if err := gui.SetKeybinding(box.Name(), gocui.MouseWheelDown, gocui.ModNone, util.ScrollDown); err != nil {
		return err
	}
	if err := gui.SetKeybinding(box.Name(), gocui.MouseWheelUp, gocui.ModNone, util.ScrollUp); err != nil {
		return err
	}
	if err := gui.SetKeybinding(box.Name(), gocui.KeyTab, gocui.ModNone, tab(box)); err != nil {
		return err
	}

	return nil
}

// SetViews sets up the views for this box, which is a simple text editor.
func (box *Script) SetViews(gui *gocui.Gui, active bool) error {
	x0, y0, x1, y1 := util.RealCoordinates(gui, &box.unit)

	if view, err := gui.SetView(box.Name(), x0, y0, x1, y1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		view.Title = box.Name()
		view.Editable = true
		view.Wrap = true
	}

	if !active {
		return nil
	}

	if _, err := gui.SetCurrentView(box.Name()); err != nil {
		return err
	}

	return nil
}

// Update updates the current command if it has changed.
func (box *Script) Update(gui *gocui.Gui, active bool) error {
	view, err := gui.View(box.Name())
	if err != nil {
		return err
	}

	x, y := view.Cursor()
	line, _ := view.Line(y)
	if x > len(line) {
		box.command = ""
	} else {
		offset := util.PositionIndex(view, x, y)
		box.command = cmds.Find(view.Buffer(), offset)
	}

	return nil
}

// Command gets the current command being worked on in the script.
func (box *Script) Command() string {
	return box.command
}

// Emulates the insertion of a tab with spaces.
func tab(box *Script) func(_ *gocui.Gui, view *gocui.View) error {
	return func(_ *gocui.Gui, view *gocui.View) error {
		for i := 0; i < box.tabSize; i++ {
			view.EditWrite(' ')
		}

		return nil
	}
}
