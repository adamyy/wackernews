package text

import "github.com/fatih/color"

type Styler func(s string, a ...interface{}) string

type Theme struct {
	PointsStyle   Styler
	RankStyle     Styler
	TitleStyle    Styler
	UrlStyle      Styler
	AuthorStyle   Styler
	TimeAgoStyle  Styler
	CommentsStyle Styler
}

var DefaultTheme = &Theme{
	PointsStyle:   color.New(color.FgYellow, color.Bold).SprintfFunc(),
	RankStyle:     color.New(color.FgYellow).SprintfFunc(),
	TitleStyle:    color.New(color.FgRed).SprintfFunc(),
	UrlStyle:      color.New(color.FgMagenta).SprintfFunc(),
	AuthorStyle:   color.New(color.FgCyan).SprintfFunc(),
	TimeAgoStyle:  color.New(color.FgGreen).SprintfFunc(),
	CommentsStyle: color.New(color.FgBlue).SprintfFunc(),
}

func (t *Theme) StylePoints(s string) string {
	return t.PointsStyle(s)
}

func (t *Theme) StyleRank(s string) string {
	return t.RankStyle(s)
}

func (t *Theme) StyleTitle(s string) string {
	return t.TitleStyle(s)
}

func (t *Theme) StyleUrl(s string) string {
	return t.UrlStyle(s)
}

func (t *Theme) StyleAuthor(s string) string {
	return t.AuthorStyle(s)
}

func (t *Theme) StyleTimeAgo(s string) string {
	return t.TimeAgoStyle(s)
}

func (t *Theme) StyleComments(s string) string {
	return t.CommentsStyle(s)
}
