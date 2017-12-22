package views

import (
	"regexp"
	"strings"

	"github.com/bryce/bashly/boxes/util"
	"github.com/jroimartin/gocui"
)

// Search holds data for a search view.
type Search struct {
	name    string
	matches [][]int
	next    int
}

const searchSuffix = " (search)"

// NewSearch creates a search view.
func NewSearch(gui *gocui.Gui, boxName string, toggle func(gui *gocui.Gui, _ *gocui.View) error) (*Search, error) {
	s := &Search{}
	s.name = boxName + searchSuffix
	if err := gui.SetKeybinding(s.name, gocui.KeyCtrlF, gocui.ModNone, toggle); err != nil {
		return nil, err
	}
	if err := gui.SetKeybinding(s.name, gocui.KeyEnter, gocui.ModNone, search(s)); err != nil {
		return nil, err
	}
	return s, nil
}

// Set sets a search view.
func (s *Search) Set(gui *gocui.Gui, x0, y0, x1, y1 int) error {
	if view, err := gui.SetView(s.name, x0, y0, x1, y1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		view.Wrap = true
		view.Editable = true
		if _, err := gui.SetCurrentView(s.name); err != nil {
			return err
		}
	}

	return nil
}

// Delete deletes a search view.
func (s *Search) Delete(gui *gocui.Gui) error {
	if err := gui.DeleteKeybinding(s.name, gocui.KeyCtrlF, gocui.ModNone); err != nil {
		return err
	}
	if err := gui.DeleteKeybinding(s.name, gocui.KeyEnter, gocui.ModNone); err != nil {
		return err
	}

	bName := strings.Replace(s.name, searchSuffix, "", 1)
	gui.DeleteKeybinding(bName, gocui.KeyEnter, gocui.ModNone)

	if err := gui.DeleteView(s.name); err != nil {
		return err
	}

	return nil
}

func (s *Search) boxView(gui *gocui.Gui) (*gocui.View, error) {
	bName := strings.Replace(s.name, searchSuffix, "", 1)
	bView, err := gui.View(bName)
	if err != nil {
		return nil, err
	}

	return bView, nil
}

func search(s *Search) func(gui *gocui.Gui, view *gocui.View) error {
	return func(gui *gocui.Gui, view *gocui.View) error {
		// Empty search
		if len(view.Buffer()) <= 1 {
			return nil
		}

		// Find matches
		bView, err := s.boxView(gui)
		if err != nil {
			return err
		}
		exp := view.Buffer()
		exp = "(?i)" + exp[:len(exp)-1]
		s.matches, err = matches(bView, exp)
		if err != nil || len(s.matches) <= 0 {
			return nil
		}

		// Setup iteration through matches
		gui.SetCurrentView(bView.Name())
		nextMatch(s)(gui, bView)
		gui.DeleteKeybinding(bView.Name(), gocui.KeyEnter, gocui.ModNone)
		if err := gui.SetKeybinding(bView.Name(), gocui.KeyEnter, gocui.ModNone, nextMatch(s)); err != nil {
			return err
		}

		return nil
	}
}

func matches(view *gocui.View, expression string) ([][]int, error) {
	buff := view.Buffer()

	re, err := regexp.Compile(expression)
	if err != nil {
		return nil, err
	}
	matches := re.FindAllStringIndex(buff, -1)

	x, y := view.Cursor()
	idx := util.PositionIndex(view, x, y)
	first := 0
	for i, match := range matches {
		if match[0] >= idx {
			first = i
			break
		}
	}
	matches = append(matches[first:], matches[:first]...)

	return matches, nil
}

func nextMatch(s *Search) func(_ *gocui.Gui, view *gocui.View) error {
	return func(_ *gocui.Gui, view *gocui.View) error {
		_, maxY := view.Size()

		x, y := util.IndexPosition(view, s.matches[s.next][0])
		oy := y - maxY/2
		if oy < 0 {
			oy = 0
		}

		view.SetOrigin(0, oy)
		y = y - oy
		view.SetCursor(x, y)

		// Prevents scrolling past the end of the manual page
		_, err := view.Line(maxY)
		for err != nil && oy != 0 {
			oy--
			view.SetOrigin(0, oy)
			y++
			view.SetCursor(x, y)
			_, err = view.Line(maxY)
		}

		s.next++
		s.next %= len(s.matches)

		return nil
	}
}
