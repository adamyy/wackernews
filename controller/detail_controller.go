package controller

import (
	"context"
	"fmt"
	"time"

	"github.com/adamyy/wackernews/view/text"

	"github.com/adamyy/wackernews/news"
	"github.com/adamyy/wackernews/view"
	"github.com/jroimartin/gocui"
)

type DetailController struct {
	id int

	detailView  *view.DetailView
	messageView *view.MessageView

	nav Navigator
}

func NewDetailController(nav Navigator, id int) *DetailController {
	c := &DetailController{
		id: id,
		messageView: view.NewMessageView(
			view.Name(fmt.Sprintf("DetailController@%d:MessageView", id)),
		),
		detailView: view.NewDetailView(
			view.Name(fmt.Sprintf("DetailController@%d:DetailView", id)),
		),
		nav: nav,
	}

	return c
}

func (dc *DetailController) Layout(g *gocui.Gui) error {
	termX, termY := g.Size()
	startX, startY := 0, 0
	endX, endY := termX-1, termY-1

	{ // setup detail view
		dv := dc.detailView
		_ = dv.Set(view.Dimension(startX, startY, endX, endY))
		v, err := g.SetView(dv.Name(), startX, startY, endX, endY)
		viewCreated := err == gocui.ErrUnknownView
		if err != nil && err != gocui.ErrUnknownView {
			return err
		}
		if viewCreated {
			if err := view.BindKeys(g, dv); err != nil {
				return err
			}
		}
		if err := dv.Draw(v); err != nil {
			return err
		}
	}

	return nil

}

type detailResponse struct {
	detail *news.Detail
	err    error
}

func (dc *DetailController) Focus(g *gocui.Gui) error {
	ticker := time.NewTicker(time.Second / 5)
	nextLoadingText := text.LoadingText()

	if err := dc.BindKeys(g); err != nil {
		return err
	}

	client, err := news.NewClient()
	if err != nil {
		return err
	}
	ctx, _ := context.WithCancel(context.Background()) // TODO cancel

	result := make(chan *detailResponse)

	go func() {
		detail, err := client.GetDetail(ctx, dc.id)
		result <- &detailResponse{detail: detail, err: err}
	}()

	go func() {
		for {
			select {
			case r, ok := <-result:
				if ok {
					g.Update(func(g *gocui.Gui) error {
						if r.err != nil {
							dc.messageView.SetMessage(r.err.Error())
						} else {
							dc.detailView.SetDetail(r.detail)
							if _, err := g.SetCurrentView(dc.detailView.Name()); err != nil {
								return err
							}
						}
						return nil
					})
				}
				ticker.Stop()
			case <-ticker.C:
				g.Update(func(g *gocui.Gui) error {
					dc.messageView.SetMessage(nextLoadingText())
					return nil
				})
			}
		}
	}()

	return nil
}

func (dc *DetailController) UnFocus(g *gocui.Gui) error {
	return nil
}

func (dc *DetailController) BindKeys(g *gocui.Gui) error {
	if err := view.BindGlobalKeys(g); err != nil {
		return err
	}

	onKeyEscape := func(g *gocui.Gui, gv *gocui.View) error {
		dc.nav.Pop()
		return nil
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlW, gocui.ModNone, onKeyEscape); err != nil {
		return err
	}

	return nil
}
