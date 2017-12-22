/*
Package boxes implements the functionality for setting up and managing
boxes. A box is an abstraction of the views and functionality of a component
in bashly.
*/
package boxes

import (
	"errors"
	"fmt"

	"github.com/jroimartin/gocui"
)

// Box abstracts a set of views and the functionality associated with them.
type Box interface {
	Name() string
	Setup(gui *gocui.Gui, boxs *Boxes) error
	SetViews(gui *gocui.Gui, active bool) error
	Update(gui *gocui.Gui, active bool) error
}

// Boxes is a collection of boxes.
type Boxes struct {
	boxes   []Box
	current int
}

// Config is the configuration format for boxes.
type Config struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	RefName string `json:"refName"`
	X0      int    `json:"x0"`
	Y0      int    `json:"y0"`
	X1      int    `json:"x1"`
	Y1      int    `json:"y1"`
	TabSize int    `json:"tabSize"`
}

// New creates the boxes for the given box configurations.
func New(cfgs []Config) (*Boxes, error) {
	boxs := &Boxes{}
	for i, cfg := range cfgs {
		var box Box
		switch cfg.Type {
		case "Script":
			box = NewScript(&cfg)
		case "Manual":
			box = NewManual(&cfg)
		default:
			return nil, fmt.Errorf("box number %d is of invalid type", i+1)
		}

		boxs.boxes = append(boxs.boxes, box)
	}

	if _, ok := boxs.boxes[0].(*Script); !ok {
		return nil, errors.New("script box must be first in configuration file")
	}

	return boxs, nil
}

// Box gets the box with the passed name.
func (boxs *Boxes) Box(name string) (Box, error) {
	for _, box := range boxs.boxes {
		if box.Name() == name {
			return box, nil
		}
	}

	return nil, errors.New("box with name " + name + " not found")
}

// Setup sets up all of the boxes.
func (boxs *Boxes) Setup(gui *gocui.Gui) error {
	for _, box := range boxs.boxes {
		if err := box.Setup(gui, boxs); err != nil {
			return err
		}
	}

	return nil
}

// SetViews sets the views for all of the boxes.
func (boxs *Boxes) SetViews(gui *gocui.Gui) error {
	for i, box := range boxs.boxes {
		if err := box.SetViews(gui, i == boxs.current); err != nil {
			return err
		}
	}

	return nil
}

// Update updates all of the boxes.
func (boxs *Boxes) Update(gui *gocui.Gui) error {
	for i, box := range boxs.boxes {
		if err := box.Update(gui, i == boxs.current); err != nil {
			return err
		}
	}

	return nil
}

// Current gets the box that is currently active.
func (boxs *Boxes) Current() Box {
	return boxs.boxes[boxs.current]
}

// Next makes the next box active.
func (boxs *Boxes) Next(gui *gocui.Gui) {
	boxs.current = (boxs.current + 1) % len(boxs.boxes)
}

// Previous makes the previous box active.
func (boxs *Boxes) Previous(gui *gocui.Gui) {
	boxs.current = boxs.current - 1
	if boxs.current < 0 {
		boxs.current = len(boxs.boxes) - 1
	}
}
