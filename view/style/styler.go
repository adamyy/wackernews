package style

type Styler interface {
	Style(s string) string
}

type Theme struct {
	PointsStyle   Styler
	RankStyle     Styler
	TitleStyle    Styler
	UrlStyle      Styler
	AuthorStyle   Styler
	TimeAgoStyle  Styler
	CommentsStyle Styler
}

var DefaultTheme = ANSITheme
