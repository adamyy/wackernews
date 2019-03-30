// A view takes some displayable data such as a list of items,
// render (format and style) the data to a given sink
//
// A view also expose a set of KeyBindings for the controller
//
package view

import (
	"github.com/adamyy/wackernews/view/style"
	"github.com/jroimartin/gocui"
)

type View interface {
	Name() string
	Draw(v *gocui.View) error
	KeyBindings() KeyBindings
}

type Prop struct {
	name string

	*dimension
	*style.Theme
}

type point struct {
	x, y int
}

type dimension struct {
	start, end *point
}

func (d dimension) Size() (int, int) {
	width := d.end.x - d.start.x
	height := d.end.y - d.start.y
	return width, height
}

type PropOption func(*Prop) error

func (p *Prop) SetProp(opts ...PropOption) error {
	for _, option := range opts {
		err := option(p)
		if err != nil {
			return err
		}
	}
	return nil
}

func Dimension(startX, startY, endX, endY int) PropOption {
	return func(prop *Prop) error {
		prop.dimension = &dimension{
			start: &point{x: startX, y: startY},
			end:   &point{x: endX, y: endY},
		}
		return nil
	}
}

func Theme(theme *style.Theme) PropOption {
	return func(prop *Prop) error {
		prop.Theme = theme
		return nil
	}
}

func Name(name string) PropOption {
	return func(prop *Prop) error {
		prop.name = name
		return nil
	}
}

var (
	DefaultPropOptions = []PropOption{
		Name("DefaultViewName"),
		Theme(style.DefaultTheme),
	}
)
