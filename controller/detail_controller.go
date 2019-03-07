package controller

import (
	"context"
	"fmt"
	"time"

	"github.com/adamyy/hackernews/view/text"

	"github.com/adamyy/hackernews/news"
	"github.com/adamyy/hackernews/view"
	"github.com/jroimartin/gocui"
)

type detailResponse struct {
	detail *news.Detail
	err    error
}

type DetailController struct {
	id int

	detailView  *view.DetailView
	messageView *view.MessageView

	reqResult chan *detailResponse
	reqCancel context.CancelFunc
}

func NewDetailController(id int) (*DetailController, error) {
	client, err := news.NewClient()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	c := &DetailController{
		id:        id,
		reqCancel: cancel,
		reqResult: make(chan *detailResponse),
		messageView: view.NewMessageView(
			view.Name(fmt.Sprintf("DetailController@%d:MessageView", id)),
		),
		detailView: view.NewDetailView(
			view.Name(fmt.Sprintf("DetailController@%d:DetailView", id)),
		),
	}

	go func() {
		detail, err := client.GetDetail(ctx, id)
		c.reqResult <- &detailResponse{detail: detail, err: err}
	}()

	return c, nil
}

func (dc *DetailController) Layout(g *gocui.Gui) error {
	termX, termY := g.Size()
	startX, startY := 0, 0
	endX, endY := termX-1, termY-1

	if err := view.BindGlobalKeys(g); err != nil {
		return err
	}

	{ // setup detail view
		dv := dc.detailView
		_ = dv.SetProp(view.Dimension(startX, startY, endX, endY))
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

func (dc *DetailController) Focus(g *gocui.Gui) error {
	ticker := time.NewTicker(time.Second / 5)
	nextLoadingText := text.LoadingText()

	for {
		select {
		case r, ok := <-dc.reqResult:
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
			return nil
		case <-ticker.C:
			g.Update(func(g *gocui.Gui) error {
				dc.messageView.SetMessage(nextLoadingText())
				return nil
			})
		}
	}
}

func (dc *DetailController) UnFocus(g *gocui.Gui) error {
	close(dc.reqResult)
	dc.reqCancel()
	return nil
}
