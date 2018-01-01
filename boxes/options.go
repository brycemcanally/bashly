package boxes

import (
	"errors"
	"reflect"

	"github.com/bryce/bashly/boxes/util"
	"github.com/bryce/bashly/manual"
	"github.com/jroimartin/gocui"
)

// Options is a type of box that holds the current options for a command.
type Options struct {
	name    string
	refName string
	unit    util.Coordinates
	script  *Script
	options []string
}

// NewOptions creates a new options box.
func NewOptions(cfg *Config) *Options {
	box := &Options{}
	box.name = cfg.Name
	box.refName = cfg.RefName
	box.unit.X0 = cfg.X0
	box.unit.Y0 = cfg.Y0
	box.unit.X1 = cfg.X1
	box.unit.Y1 = cfg.Y1

	return box
}

// Name returns the name associated with this box.
func (box *Options) Name() string {
	return box.name
}

// Setup sets up the reference box and keybindings for this box.
func (box *Options) Setup(gui *gocui.Gui, boxs *Boxes) error {
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

	return gui.SetKeybinding(box.Name(), gocui.MouseWheelUp, gocui.ModNone, util.ScrollUp)
}

// SetViews sets up the views for this box.
func (box *Options) SetViews(gui *gocui.Gui, active bool) error {
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

// Update updates the options if they have changed.
func (box *Options) Update(gui *gocui.Gui, active bool) error {
	view, err := gui.View(box.Name())
	if err != nil {
		return err
	}

	cmd, err := box.script.Command()
	if err != nil {
		view.Clear()
		box.options = []string{}
		return nil
	}

	if reflect.DeepEqual(cmd.Options, box.options) {
		return nil
	}

	view.Clear()
	view.SetOrigin(0, 0)
	view.SetCursor(0, 0)

	box.options = cmd.Options
	maxX, _ := view.Size()
	page, err := manual.GetOptions(cmd, maxX)
	if err == nil {
		view.Write(page)
	}

	return nil
}
