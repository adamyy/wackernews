package view

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/net/html"

	"github.com/adamyy/hackernews/news"
	"github.com/adamyy/hackernews/view/text"
	"github.com/jroimartin/gocui"
)

type DetailView struct {
	*Prop

	detail *news.Detail

	tainted bool
	mutex   sync.Mutex
}

func NewDetailView(opts ...PropOption) *DetailView {
	v := &DetailView{Prop: &Prop{}}
	_ = v.SetProp(DefaultPropOptions...)
	_ = v.SetProp(opts...) // TODO refactor
	return v
}

func (dv *DetailView) Draw(v *gocui.View) error {
	dv.mutex.Lock()
	defer dv.mutex.Unlock()

	if !dv.tainted { // no need re-rendered
		return nil
	}
	defer func() { dv.tainted = false }()

	v.Clear()
	v.Frame = true
	v.Title = dv.Name()

	for _, line := range dv.render() {
		if _, err := fmt.Fprintln(v, line); err != nil {
			return err
		}
	}

	return nil
}

func (dv *DetailView) KeyBindings() KeyBindings {
	onKeyArrowUp := func(g *gocui.Gui, gv *gocui.View) error {
		_ = ScrollLines(gv, -1)
		return nil
	}

	onKeyArrowDown := func(g *gocui.Gui, gv *gocui.View) error {
		_ = ScrollLines(gv, 1)
		return nil
	}

	onKeyArrowLeft := func(g *gocui.Gui, gv *gocui.View) error {
		_ = ScrollPage(gv, -1)
		return nil
	}

	onKeyArrowRight := func(g *gocui.Gui, gv *gocui.View) error {
		_ = ScrollPage(gv, 1)
		return nil
	}

	return KeyBindings{
		gocui.ModNone: {
			gocui.KeyArrowUp:    onKeyArrowUp,
			gocui.KeyArrowDown:  onKeyArrowDown,
			gocui.KeyArrowLeft:  onKeyArrowLeft,
			gocui.KeyArrowRight: onKeyArrowRight,
		},
	}
}

func (dv *DetailView) SetDetail(detail *news.Detail) {
	dv.mutex.Lock()
	defer dv.mutex.Unlock()

	dv.tainted = true
	dv.detail = detail
}

func (dv *DetailView) Name() string {
	return dv.name
}

func (dv *DetailView) render() []string {
	var lines []string
	lines = append(lines, dv.renderHeader()...)
	lines = append(lines, dv.renderContent()...)
	lines = append(lines, dv.renderComments()...)
	return lines
}

// title (url)
// points by author time-ago | comments
func (dv *DetailView) renderHeader() []string {
	d := dv.detail
	s := dv.theme

	title := s.TitleStyle.Style(d.Title)
	url := s.UrlStyle.Style(d.Url)
	points := s.PointsStyle.Style(strconv.Itoa(d.Points))
	author := s.AuthorStyle.Style(d.User)
	timeAgo := s.TimeAgoStyle.Style(d.TimeAgo)
	comments := s.CommentsStyle.Style(strconv.Itoa(d.CommentsCount))

	return []string{
		fmt.Sprintf("%s (%s)", title, url),
		fmt.Sprintf("%s points by %s %s | %s comments", points, author, timeAgo, comments),
	}
}

func (dv *DetailView) renderContent() []string {
	d := dv.detail
	width, _ := dv.Size()
	return text.Justify(d.Content, width-1, false)
}

func (dv *DetailView) renderComments() []string {
	d := dv.detail
	var lines []string
	for _, comment := range d.Comments {
		rendered := dv.renderComment(comment)
		lines = append(lines, rendered...)
	}
	return lines
}

// ▲ user time ago
// comment text
func (dv *DetailView) renderComment(comment *news.Detail) []string {
	if comment.Dead || comment.Deleted {
		return nil
	}
	s := dv.theme
	indent := strings.Repeat("\t\t\t", comment.Level)
	width, _ := dv.Size()
	renderWidth := width - 1 - len(indent)

	var lines []string
	user := s.AuthorStyle.Style(comment.User)
	timeAgo := s.TimeAgoStyle.Style(comment.TimeAgo)
	header := fmt.Sprintf("%s▲ %s %s", indent, user, timeAgo)
	lines = append(lines, header)

	commentLines := strings.Split(html.UnescapeString(comment.Content), "<p>")
	for _, commentLine := range commentLines {
		justified := text.Justify(commentLine, renderWidth, false)
		for _, justifiedLine := range justified {
			lines = append(lines, indent+justifiedLine)
		}
	}

	for _, subComment := range comment.Comments {
		rendered := dv.renderComment(subComment)
		lines = append(lines, rendered...)
	}
	return lines
}
