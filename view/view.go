package view

import (
	"github.com/adamyy/wackernews/view/text"
	"github.com/jroimartin/gocui"
)

// View takes some displayable data such as a list of items
// render (format and style) the data to a given source
// A view also expose a set of KeyBindings for the controller
type View interface {
	Name() string
	Size() (int, int)
	Draw(v *gocui.View) error
	KeyBindings() KeyBindings
}

type Prop struct {
	name string

	*dimension
	*text.Theme
}

type PropOption func(*Prop) error

// DefaultProp returns a prop instance with default options applied
func DefaultProp() *Prop {
	p := &Prop{}
	p.Set(Name("DefaultViewName"), Theme(text.DefaultTheme))
	return p
}

type dimension struct {
	start, end *point
}

type point struct {
	x, y int
}

func (d dimension) Size() (int, int) {
	width := d.end.x - d.start.x
	height := d.end.y - d.start.y
	return width, height
}

// Set applies the given prop changes
func (p *Prop) Set(opts ...PropOption) error {
	for _, option := range opts {
		err := option(p)
		if err != nil {
			return err
		}
	}
	return nil
}

// Dimension sets a view's start and end location
func Dimension(startX, startY, endX, endY int) PropOption {
	return func(p *Prop) error {
		p.dimension = &dimension{
			start: &point{x: startX, y: startY},
			end:   &point{x: endX, y: endY},
		}
		return nil
	}
}

// Theme sets a view's text styles
func Theme(theme *text.Theme) PropOption {
	return func(p *Prop) error {
		p.Theme = theme
		return nil
	}
}

// Name sets a view's unique identifier
func Name(name string) PropOption {
	return func(p *Prop) error {
		p.name = name
		return nil
	}
}
