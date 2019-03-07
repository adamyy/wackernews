package view

import "github.com/jroimartin/gocui"

func ScrollLines(v *gocui.View, lines int) error {
	contentLength := len(v.BufferLines())
	_, windowY := v.Size()
	startX, startY := v.Origin()
	endX, endY := startX, startY+lines

	if endY+windowY > contentLength { // go beyond bottom
		return ScrollToBottom(v)
	}
	if endY < 0 { // go above top
		return ScrollToTop(v)
	}

	return v.SetOrigin(endX, endY)
}

func ScrollPage(v *gocui.View, pages int) error {
	_, windowY := v.Size()
	return ScrollLines(v, windowY*pages)
}

func ScrollToTop(v *gocui.View) error {
	startX, _ := v.Origin()
	endX, endY := startX, 0

	return v.SetOrigin(endX, endY)
}

func ScrollToBottom(v *gocui.View) error {
	contentLength := len(v.BufferLines())
	_, windowY := v.Size()
	startX, _ := v.Origin()
	endX, endY := startX, contentLength-windowY

	if windowY >= contentLength { // if the view can fit all content
		return nil // no-op
	}

	return v.SetOrigin(endX, endY)
}
