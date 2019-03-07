package controller

import (
	"context"
	"fmt"
	"time"

	"github.com/adamyy/hackernews/news"
	"github.com/adamyy/hackernews/view"
	"github.com/adamyy/hackernews/view/text"
	"github.com/jroimartin/gocui"
)

type feedResponse struct {
	feed *news.Feed
	err  error
}

type FeedController struct {
	kind news.FeedKind
	page int

	feedView    *view.FeedView
	messageView *view.MessageView

	reqResult chan *feedResponse
	reqCancel context.CancelFunc
}

func NewFeedController(kind news.FeedKind, page int) (*FeedController, error) {
	client, err := news.NewClient()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	c := &FeedController{
		kind:      kind,
		page:      page,
		reqCancel: cancel,
		reqResult: make(chan *feedResponse),
		messageView: view.NewMessageView(
			view.Name(fmt.Sprintf("FeedController@%d/%d:MessageView", kind, page)),
		),
		feedView: view.NewFeedView(
			view.Name(fmt.Sprintf("FeedController@%d/%d:FeedView", kind, page)),
		),
	}

	go func() {
		feed, err := client.GetFeed(ctx, kind, page)
		c.reqResult <- &feedResponse{feed: feed, err: err}
	}()

	return c, nil
}

func (fc *FeedController) Layout(g *gocui.Gui) error {
	termX, termY := g.Size()
	startX, startY := 0, 0
	endX, endY := termX-1, termY-1

	if err := view.BindGlobalKeys(g); err != nil {
		return err
	}

	{ // setup feed view
		fv := fc.feedView
		_ = fv.SetProp(view.Dimension(startX, startY, endX, endY*10/11))
		v, err := g.SetView(fc.feedView.Name(), startX, startY, endX, endY*10/11)
		viewCreated := err == gocui.ErrUnknownView
		if err != nil && err != gocui.ErrUnknownView {
			return err
		}
		if viewCreated {
			if err := view.BindKeys(g, fv); err != nil {
				return err
			}
		}
		if err := fv.Draw(v); err != nil {
			return err
		}
	}

	{ // setup message view
		mv := fc.messageView
		_ = mv.SetProp(view.Dimension(startX, endY*11/12, endX, endY))
		v, err := g.SetView(fc.messageView.Name(), startX, endY*11/12, endX, endY)
		viewCreated := err == gocui.ErrUnknownView
		if err != nil && err != gocui.ErrUnknownView {
			return err
		}
		if viewCreated {
			if err := view.BindKeys(g, mv); err != nil {
				return err
			}
		}
		if err := mv.Draw(v); err != nil {
			return err
		}
	}

	return nil
}

func (fc *FeedController) Focus(g *gocui.Gui) error {
	ticker := time.NewTicker(time.Second / 5)
	nextLoadingText := text.LoadingText()

	for {
		select {
		case r, ok := <-fc.reqResult:
			if ok {
				g.Update(func(g *gocui.Gui) error {
					if r.err != nil {
						fc.messageView.SetMessage(r.err.Error())
					} else {
						fc.feedView.SetFeed(r.feed)
						if _, err := g.SetCurrentView(fc.feedView.Name()); err != nil {
							return err
						}
					}
					return nil
				})
			}
			return nil
		case <-ticker.C:
			g.Update(func(g *gocui.Gui) error {
				fc.messageView.SetMessage(nextLoadingText())
				return nil
			})
		}
	}
}

func (fc *FeedController) UnFocus(g *gocui.Gui) error {
	close(fc.reqResult)
	fc.reqCancel()
	return nil
}
