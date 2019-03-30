package app

import (
	"log"

	"github.com/adamyy/wackernews/controller"
	"github.com/adamyy/wackernews/news"
	"github.com/jroimartin/gocui"
)

type App struct {
	gui   *gocui.Gui
	stack []controller.Controller
}

func NewApp() (*App, error) {
	gui, err := gocui.NewGui(gocui.Output256)
	if err != nil {
		return nil, err
	}
	return &App{gui: gui}, nil
}

func (a *App) Init() error {
	fc := controller.NewFeedController(a, news.KindNews, 1)
	a.Push(fc)

	if err := a.gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		return err
	}
	return nil
}

func (a *App) Push(c controller.Controller) {
	if l := len(a.stack); l > 0 {
		curr := a.stack[l-1]
		if err := curr.UnFocus(a.gui); err != nil {
			log.Panicln(err)
		}
	}

	a.stack = append(a.stack, c)
	a.gui.SetManager(c)

	if err := c.Focus(a.gui); err != nil {
		log.Panicln(err)
	}
}

func (a *App) Pop() {
	if l := len(a.stack); l > 0 {
		curr := a.stack[l-1]
		if err := curr.UnFocus(a.gui); err != nil {
			log.Panicln(err)
		}
		a.stack = a.stack[:l-1]
	}

	if l := len(a.stack); l > 0 {
		prev := a.stack[l-1]
		a.gui.SetManager(prev)
		if err := prev.Focus(a.gui); err != nil {
			log.Panicln(err)
		}
	}
}

func (a *App) Close() {
	a.gui.Close()
}
