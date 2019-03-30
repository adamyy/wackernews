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

func (t *Theme) StylePoints(s string) string {
	return t.PointsStyle.Style(s)
}

func (t *Theme) StyleRank(s string) string {
	return t.RankStyle.Style(s)
}

func (t *Theme) StyleTitle(s string) string {
	return t.TitleStyle.Style(s)
}

func (t *Theme) StyleUrl(s string) string {
	return t.UrlStyle.Style(s)
}

func (t *Theme) StyleAuthor(s string) string {
	return t.AuthorStyle.Style(s)
}

func (t *Theme) StyleTimeAgo(s string) string {
	return t.TimeAgoStyle.Style(s)
}

func (t *Theme) StyleComments(s string) string {
	return t.CommentsStyle.Style(s)
}
