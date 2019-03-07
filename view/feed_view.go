package view

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/adamyy/hackernews/news"
	"github.com/jroimartin/gocui"
)

type FeedView struct {
	*Prop

	index int
	feed  *news.Feed

	tainted bool
	mutex   sync.Mutex
}

func NewFeedView(opts ...PropOption) *FeedView {
	v := &FeedView{Prop: &Prop{}}
	_ = v.SetProp(DefaultPropOptions...)
	_ = v.SetProp(opts...) // TODO refactor
	return v
}

func (fv *FeedView) Draw(v *gocui.View) error {
	fv.mutex.Lock()
	defer fv.mutex.Unlock()

	if fv.feed == nil { // content is missing, skipping rendering
		return nil
	}

	if !fv.tainted { // no need to re-render
		return nil
	}
	defer func() { fv.tainted = false }()

	v.Clear()
	v.Frame = true
	v.Title = fv.feed.Kind.ReadableString()
	v.Highlight = true
	v.SelBgColor = gocui.ColorBlack

	lines := fv.Render()
	_, err := fmt.Fprintf(v, strings.Join(lines, "\n"))

	return err
}

func (fv *FeedView) SetFeed(feed *news.Feed) {
	fv.mutex.Lock()
	defer fv.mutex.Unlock()

	fv.feed = feed
	fv.index = 0
	fv.tainted = true
}

func (fv *FeedView) Name() string {
	return fv.name
}

func (fv *FeedView) KeyBindings() KeyBindings {
	onKeyArrowUp := func(g *gocui.Gui, gv *gocui.View) error {
		_ = MoveCursor(gv, -2)
		if fv.index > 0 {
			fv.index = fv.index - 1
		}
		return nil
	}

	onKeyArrowDown := func(g *gocui.Gui, gv *gocui.View) error {
		_ = MoveCursor(gv, 2)
		if fv.index+1 < len(fv.feed.Items) {
			fv.index = fv.index + 1
		}
		return nil
	}

	onKeyArrowLeft := func(g *gocui.Gui, gv *gocui.View) error {
		_ = MoveCursorTop(gv)
		fv.index = 0
		return nil
	}

	onKeyArrowRight := func(g *gocui.Gui, gv *gocui.View) error {
		_ = MoveCursorBottom(gv)
		_ = MoveCursor(gv, -1)
		fv.index = len(fv.feed.Items) - 1
		return nil
	}

	onKeyEnter := func(g *gocui.Gui, gv *gocui.View) error {
		return nil
	}

	return KeyBindings{
		gocui.ModNone: {
			gocui.KeyArrowUp:    onKeyArrowUp,
			gocui.KeyArrowDown:  onKeyArrowDown,
			gocui.KeyArrowLeft:  onKeyArrowLeft,
			gocui.KeyArrowRight: onKeyArrowRight,
			gocui.KeyEnter:      onKeyEnter,
		},
	}
}

func (fv *FeedView) SelectedItem() *news.Item {
	if fv.feed == nil {
		return nil
	}
	return fv.feed.Items[fv.index]
}

// rank    	title (url)
//			points by author time-ago | comments
func (fv *FeedView) Render() []string {
	feed := fv.feed
	lines := make([]string, len(feed.Items))
	maxRank := feed.Page * len(feed.Items)
	indent := strings.Repeat(" ", len(strconv.Itoa(maxRank))+3)
	s := fv.theme
	for index, item := range feed.Items {
		rankStr := strconv.Itoa(index + 1)

		rank := s.RankStyle.Style(fmt.Sprintf("[%s]%s", rankStr, indent[len(rankStr)+3:]))
		title := s.TitleStyle.Style(item.Title)
		url := s.UrlStyle.Style(item.Url)
		points := s.PointsStyle.Style(strconv.Itoa(item.Points))
		author := s.AuthorStyle.Style(item.User)
		timeAgo := s.TimeAgoStyle.Style(item.TimeAgo)
		comments := s.CommentsStyle.Style(strconv.Itoa(item.CommentsCount))

		firstLine := fmt.Sprintf("%s %s (%s)", rank, title, url)
		secondLine := fmt.Sprintf("%s%s points by %s %s| %s comments", indent, points, author, timeAgo, comments)

		lines[index] = firstLine + "\n" + secondLine
	}
	return lines
}
