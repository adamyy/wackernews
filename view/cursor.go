package view

import "github.com/jroimartin/gocui"

func MoveCursor(v *gocui.View, lines int) error {
	v.MoveCursor(0, lines, false)
	return nil
}

func MoveCursorTop(v *gocui.View) error {
	cursorX, _ := v.Cursor()
	return v.SetCursor(cursorX, 0)
}

func MoveCursorBottom(v *gocui.View) error {
	contentLength := len(v.BufferLines())
	cursorX, _ := v.Cursor()
	return v.SetCursor(cursorX, contentLength-1)
}
