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

type FeedController struct {
	kind news.FeedKind
	page int

	feedView    *view.FeedView
	messageView *view.MessageView

	nav Navigator
}

func NewFeedController(nav Navigator, kind news.FeedKind, page int) *FeedController {
	c := &FeedController{
		kind: kind,
		page: page,
		nav:  nav,
		feedView: view.NewFeedView(
			view.Name(fmt.Sprintf("FeedController@%d/%d:FeedView", kind, page)),
		),
		messageView: view.NewMessageView(
			view.Name(fmt.Sprintf("FeedController@%d/%d:MessageView", kind, page)),
		),
	}
	return c
}

func (fc *FeedController) Layout(g *gocui.Gui) error {
	termX, termY := g.Size()
	startX, startY := 0, 0
	endX, endY := termX-1, termY-1

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

type feedResponse struct {
	feed *news.Feed
	err  error
}

func (fc *FeedController) Focus(g *gocui.Gui) error {
	ticker := time.NewTicker(time.Second / 5)
	nextLoadingText := text.LoadingText()

	if err := fc.BindKeys(g); err != nil {
		return err
	}

	client, err := news.NewClient()
	if err != nil {
		return err
	}
	ctx, _ := context.WithCancel(context.Background()) // TODO cancel

	result := make(chan *feedResponse)

	go func() {
		feed, err := client.GetFeed(ctx, fc.kind, fc.page)
		result <- &feedResponse{feed: feed, err: err}
	}()

	go func() {
		for {
			select {
			case r, ok := <-result:
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
				ticker.Stop()
			case <-ticker.C:
				g.Update(func(g *gocui.Gui) error {
					fc.messageView.SetMessage(nextLoadingText())
					return nil
				})
			}
		}
	}()

	return nil
}

func (fc *FeedController) UnFocus(g *gocui.Gui) error {
	return nil
}

func (fc *FeedController) BindKeys(g *gocui.Gui) error {
	if err := view.BindGlobalKeys(g); err != nil {
		return err
	}

	onKeyEnter := func(g *gocui.Gui, gv *gocui.View) error {
		item := fc.feedView.SelectedItem()
		fc.messageView.SetMessage(item.Title)
		dc := NewDetailController(fc.nav, item.Id)
		fc.nav.Push(dc)
		return nil
	}

	if err := g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, onKeyEnter); err != nil {
		return err
	}

	onKeyTab := func(g *gocui.Gui, gv *gocui.View) error {
		c := NewFeedController(fc.nav, fc.kind, fc.page+1)
		fc.nav.Push(c)
		return nil
	}

	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, onKeyTab); err != nil {
		return err
	}

	return nil
}
