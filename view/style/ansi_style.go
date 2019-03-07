package style

import (
	"github.com/fatih/color"
)

type AnsiStyler struct {
	*color.Color
}

func (s *AnsiStyler) Style(str string) string {
	return s.Sprintf(str)
}

var (
	ANSITheme = &Theme{
		PointsStyle:   &AnsiStyler{Color: color.New(color.FgYellow, color.Bold)},
		RankStyle:     &AnsiStyler{Color: color.New(color.FgYellow)},
		TitleStyle:    &AnsiStyler{Color: color.New(color.FgRed)},
		UrlStyle:      &AnsiStyler{Color: color.New(color.FgMagenta)},
		AuthorStyle:   &AnsiStyler{Color: color.New(color.FgCyan)},
		TimeAgoStyle:  &AnsiStyler{Color: color.New(color.FgGreen)},
		CommentsStyle: &AnsiStyler{Color: color.New(color.FgBlue)},
	}
)
