package view

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/net/html"

	"github.com/adamyy/wackernews/news"
	"github.com/adamyy/wackernews/view/text"
	"github.com/jroimartin/gocui"
)

type DetailView struct {
	*Prop

	detail *news.Detail

	tainted bool
	mutex   sync.Mutex
}

func NewDetailView(opts ...PropOption) *DetailView {
	v := &DetailView{Prop: DefaultProp()}
	_ = v.Set(opts...)
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
	v.Title = dv.detail.Title

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

	title := dv.StyleTitle(d.Title)
	url := dv.StyleUrl(d.Url)
	points := dv.StylePoints(strconv.Itoa(d.Points))
	author := dv.StyleAuthor(d.User)
	timeAgo := dv.StyleTimeAgo(d.TimeAgo)
	comments := dv.StyleComments(strconv.Itoa(d.CommentsCount))

	return []string{
		fmt.Sprintf("%s (%s)", title, url),
		fmt.Sprintf("%s points by %s %s | %s comments", points, author, timeAgo, comments),
	}
}

func (dv *DetailView) renderContent() []string {
	d := dv.detail
	width, _ := dv.Size()
	return dv.formatContent(d.Content, width-1)
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
	indent := strings.Repeat("\t\t\t", comment.Level)
	width, _ := dv.Size()
	renderWidth := width - 2 - len(indent)

	var lines []string
	user := dv.StyleAuthor(comment.User)
	timeAgo := dv.StyleTimeAgo(comment.TimeAgo)
	header := fmt.Sprintf("%s▲ %s %s", indent, user, timeAgo)
	lines = append(lines, header)

	content := dv.formatContent(comment.Content, renderWidth)
	for _, line := range content {
		lines = append(lines, indent+"\t\t"+line)
	}

	for _, subComment := range comment.Comments {
		rendered := dv.renderComment(subComment)
		lines = append(lines, rendered...)
	}
	return lines
}

func (dv *DetailView) formatContent(content string, width int) []string {
	var lines []string
	unescaped := html.UnescapeString(content)
	split := strings.Split(unescaped, "<p>")

	for _, l := range split {
		humanized := text.Humanize(l)
		justified := text.Justify(humanized, width, false)
		lines = append(lines, justified...)
	}

	return lines
}
