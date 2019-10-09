// A controller is a virtual representation of a "page".
//
// A controller can manage multiple views, including how the views
// are positioned, which view is currently in focused,
// what data should get passed to the view for rendering, etc.
//
// A controller is also responsible for fetching data from data source.
//
package controller

import "github.com/jroimartin/gocui"

type Controller interface {
	gocui.Manager
	Focus(g *gocui.Gui) error
	UnFocus(g *gocui.Gui) error
}

type Navigator interface {
	Push(c Controller)
	Pop()
}
