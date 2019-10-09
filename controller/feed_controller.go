package controller

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/adamyy/wackernews/news"
	"github.com/adamyy/wackernews/view"
	"github.com/adamyy/wackernews/view/text"
	"github.com/jroimartin/gocui"
)

type FeedController struct {
	kind news.FeedKind
	page int

	feedView    *view.FeedView
	messageView *view.MessageView

	client *news.Client

	events     chan feedEvent
	interrupts chan interface{}

	nav Navigator
}

func NewFeedController(nav Navigator, kind news.FeedKind, page int) *FeedController {
	client, _ := news.NewClient()

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
		client:     client,
		events:     make(chan feedEvent),
		interrupts: make(chan interface{}),
	}

	return c
}

func (fc *FeedController) Layout(g *gocui.Gui) error {
	termX, termY := g.Size()
	startX, startY := 0, 0
	endX, endY := termX-1, termY-1

	{ // setup feed view
		fv := fc.feedView
		_ = fv.Set(view.Dimension(startX, startY, endX, endY))
		v, err := g.SetView(fc.feedView.Name(), startX, startY, endX, endY)
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
		midX := (startX + endX) / 2
		midY := (startX + endY) / 2
		_ = mv.Set(view.Dimension(midX-10, midY, midX+10, midY+2))
		v, err := g.SetView(fc.messageView.Name(), midX-10, midY, midX+10, midY+2)
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

type feedEvent func(fc *FeedController) error

func (fc *FeedController) listen() {
	for {
		select {
		case e := <-fc.events:
			if err := e(fc); err != nil {
				panic(err)
			}
		case <-fc.interrupts:
			return
		}
	}
}

func (fc *FeedController) Focus(g *gocui.Gui) error {
	if err := fc.BindKeys(g); err != nil {
		return err
	}

	go fc.listen()
	go fc.load(g)

	return nil
}

func (fc *FeedController) load(g *gocui.Gui) {
	done := make(chan bool)
	defer close(done)

	nextLoadingText := text.LoadingText()
	go ticks(done, func() {
		fc.events <- func(fc *FeedController) error {
			g.Update(func(g *gocui.Gui) error {
				fc.messageView.SetMessage(nextLoadingText())
				return nil
			})
			return nil
		}
	})
	if _, err := g.SetViewOnTop(fc.messageView.Name()); err != nil {
		log.Println("error setting message view on top")
	}

	ctx := context.Background()
	feed, err := fc.client.GetFeed(ctx, fc.kind, fc.page)
	done <- true

	fc.events <- func(fc *FeedController) error {
		g.Update(func(g *gocui.Gui) error {
			if err != nil {
				fc.messageView.SetMessage(err.Error())
			} else {
				fc.feedView.SetFeed(feed)
				if _, err := g.SetCurrentView(fc.feedView.Name()); err != nil {
					return err
				}
				if _, err := g.SetViewOnTop(fc.feedView.Name()); err != nil {
					return err
				}
			}
			return nil
		})
		return nil
	}
}

func (fc *FeedController) UnFocus(g *gocui.Gui) error {
	fc.interrupts <- true
	return nil
}

func (fc *FeedController) BindKeys(g *gocui.Gui) error {
	if err := view.BindGlobalKeys(g); err != nil {
		return err
	}

	{
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
	}

	{
		onKeyArrowLeft := func(g *gocui.Gui, gv *gocui.View) error {
			if fc.page > 1 {
				c := NewFeedController(fc.nav, fc.kind, fc.page-1)
				fc.nav.Push(c)
			}
			return nil
		}
		if err := g.SetKeybinding("", gocui.KeyArrowLeft, gocui.ModNone, onKeyArrowLeft); err != nil {
			return err
		}
	}

	{
		onKeyArrowRight := func(g *gocui.Gui, gv *gocui.View) error {
			c := NewFeedController(fc.nav, fc.kind, fc.page+1)
			fc.nav.Push(c)
			return nil
		}
		if err := g.SetKeybinding("", gocui.KeyArrowRight, gocui.ModNone, onKeyArrowRight); err != nil {
			return err
		}
	}

	return nil
}

func ticks(done <-chan bool, onTick func()) {
	ticker := time.NewTicker(time.Second / 5)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			onTick()
		case <-done:
			return
		}
	}
}
