package boxes

import (
	"errors"

	"github.com/bryce/bashly/boxes/util"
	"github.com/bryce/bashly/boxes/views"
	"github.com/bryce/bashly/manual"
	"github.com/jroimartin/gocui"
)

// Manual is a type of box that holds the manual pages
// associated with a script.
type Manual struct {
	name    string
	refName string
	unit    util.Coordinates
	script  *Script
	command string
	subView views.View // change name
}

// NewManual creates a new manual page box.
func NewManual(cfg *Config) *Manual {
	box := &Manual{}
	box.name = cfg.Name
	box.refName = cfg.RefName
	box.unit.X0 = cfg.X0
	box.unit.Y0 = cfg.Y0
	box.unit.X1 = cfg.X1
	box.unit.Y1 = cfg.Y1

	return box
}

// Name returns the name associated with this box.
func (box *Manual) Name() string {
	return box.name
}

// Setup sets up the reference box and keybindings for this box.
func (box *Manual) Setup(gui *gocui.Gui, boxs *Boxes) error {
	sBox, err := boxs.Box(box.refName)
	if err != nil {
		return err
	}
	var ok bool
	if box.script, ok = sBox.(*Script); !ok {
		return errors.New("reference box has wrong type")
	}

	if err := gui.SetKeybinding(box.Name(), gocui.MouseWheelDown, gocui.ModNone, util.ScrollDown); err != nil {
		return err
	}
	if err := gui.SetKeybinding(box.Name(), gocui.MouseWheelUp, gocui.ModNone, util.ScrollUp); err != nil {
		return err
	}

	return gui.SetKeybinding(box.Name(), gocui.KeyCtrlF, gocui.ModNone, search(box))
}

// SetViews sets up the views for the manual pages in this box.
func (box *Manual) SetViews(gui *gocui.Gui, active bool) error {
	x0, y0, x1, y1 := util.RealCoordinates(gui, &box.unit)

	if view, err := gui.SetView(box.Name(), x0, y0, x1, y1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		view.Title = box.Name()
		view.Wrap = true
	}

	if !active {
		if box.subView != nil {
			return toggleOff(box)(gui, nil)
		}

		return nil
	}

	if box.subView != nil {
		box.subView.Set(gui, x1-((x1-x0)/3), y0+1, x1-2, y0+3)
	} else if _, err := gui.SetCurrentView(box.Name()); err != nil {
		return err
	}

	return nil
}

// Update updates the manual pages if the script comands change.
func (box *Manual) Update(gui *gocui.Gui, active bool) error {
	view, err := gui.View(box.Name())
	if err != nil {
		return err
	}

	cmd, err := box.script.Command()
	if err != nil {
		view.Clear()
		box.command = ""
		return nil
	}

	if cmd.Name == box.command {
		return nil
	}

	view.Clear()
	view.SetOrigin(0, 0)
	view.SetCursor(0, 0)

	box.command = cmd.Name
	maxX, _ := view.Size()
	page, err := manual.Get(cmd, maxX)
	if err == nil {
		view.Write(page)
	}

	return nil
}

func toggleOff(box *Manual) func(gui *gocui.Gui, _ *gocui.View) error {
	return func(gui *gocui.Gui, _ *gocui.View) error {
		err := box.subView.Delete(gui)
		if err != nil {
			return err
		}
		box.subView = nil

		return nil
	}
}

func search(box *Manual) func(gui *gocui.Gui, _ *gocui.View) error {
	return func(gui *gocui.Gui, _ *gocui.View) error {
		toggleFunc := toggleOff(box)

		if box.subView != nil {
			return toggleFunc(gui, nil)
		}

		s, err := views.NewSearch(gui, box.Name(), toggleFunc)
		if err != nil {
			return err
		}
		box.subView = s

		return nil
	}
}
