package main

import (
	"github.com/adamyy/hackernews/controller"
	"github.com/adamyy/hackernews/news"
	"github.com/jroimartin/gocui"
)

type App struct {
	gui  *gocui.Gui
	flow []controller.Controller
}

func NewApp() (*App, error) {
	gui, err := gocui.NewGui(gocui.Output256)
	if err != nil {
		return nil, err
	}
	return &App{gui: gui}, nil
}

func (app *App) Init() error {
	fc, err := controller.NewFeedController(news.KindNews, 1)
	if err != nil {
		return err
	}
	app.gui.SetManager(fc)

	app.flow = append(app.flow, fc)
	go fc.Focus(app.gui)

	if err := app.gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		return err
	}
	return nil
}

func (app *App) Close() {
	app.gui.Close()
}
