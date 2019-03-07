package view

import (
	"github.com/jroimartin/gocui"
)

type KeyBindingFunc func(*gocui.Gui, *gocui.View) error

type KeyBindings map[gocui.Modifier]map[gocui.Key]KeyBindingFunc

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func BindGlobalKeys(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	return nil
}

func BindKeys(g *gocui.Gui, v View) error {
	for mod, keyToFunc := range v.KeyBindings() {
		for key, bindingFunc := range keyToFunc {
			if err := g.SetKeybinding(v.Name(), key, mod, bindingFunc); err != nil {
				return err
			}
		}
	}

	return nil
}
