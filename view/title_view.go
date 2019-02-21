package view

import "github.com/fatih/color"

type TitleView struct {
	title   string
	loading chan bool
}

func (v *TitleView) render() string {
	select {
	case <-v.loading:
		return color.YellowString("Loading...")
	default:
		return color.WhiteString(v.title)
	}
}
