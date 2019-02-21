package view

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/adamyy/hackernews/feed"
	"github.com/fatih/color"
)

type FeedView struct {
	feed.Type
	Items  []feed.Item
	Page   int
	Cursor int
}

func NewFeedView(feedType feed.Type, items []feed.Item, page int) *FeedView {
	return &FeedView{Type: feedType, Items: items, Page: page, Cursor: 0}
}

func (v *FeedView) Render() string {
	header := v.renderHeader()
	items := v.renderItems()
	return strings.Join([]string{header, items}, "\n")
}

var headerItemColor = color.New(color.BgRed)
var headerItemColorSelected = color.New(color.BgHiRed).Add(color.Bold)

func (v *FeedView) renderHeader() string {
	headerItems := make([]string, len(feed.AllTypes))
	for i, t := range feed.AllTypes {
		if t == v.Type {
			headerItems[i] = headerItemColorSelected.Sprintf("%s    ", t.ReadableString())
		} else {
			headerItems[i] = headerItemColor.Sprintf("%s    ", t.ReadableString())
		}
	}
	return strings.Join(headerItems, "")
}

type ItemColorSet struct {
	PointsColor   *color.Color
	RankColor     *color.Color
	TitleColor    *color.Color
	UrlColor      *color.Color
	AuthorColor   *color.Color
	TimeAgoColor  *color.Color
	CommentsColor *color.Color
}

var (
	defaultItemColorSet = &ItemColorSet{
		PointsColor:   color.New(color.FgYellow).Add(color.Bold),
		RankColor:     color.New(color.FgHiYellow),
		TitleColor:    color.New(color.FgWhite),
		UrlColor:      color.New(color.FgWhite),
		AuthorColor:   color.New(color.FgCyan),
		TimeAgoColor:  color.New(color.FgGreen),
		CommentsColor: color.New(color.FgBlue),
	}
	selectedItemColorSet = &ItemColorSet{
		PointsColor:   color.New(color.FgYellow).Add(color.Bold),
		RankColor:     color.New(color.FgHiYellow),
		TitleColor:    color.New(color.FgWhite).Add(color.Bold),
		UrlColor:      color.New(color.FgWhite),
		AuthorColor:   color.New(color.FgCyan),
		TimeAgoColor:  color.New(color.FgGreen),
		CommentsColor: color.New(color.FgBlue),
	}
)

// rank    	title (url)
//			points by author time-ago | comments
func (v *FeedView) renderItems() string {
	lines := make([]string, len(v.Items))
	for index, item := range v.Items {
		cs := defaultItemColorSet
		if index == v.Cursor {
			cs = selectedItemColorSet
		}

		rank := cs.RankColor.Sprintf(strconv.Itoa((v.Page-1)*30 + index + 1))
		title := cs.TitleColor.Sprintf(item.Title)
		url := cs.UrlColor.Sprintf(item.Url)
		points := cs.PointsColor.Sprintf(strconv.Itoa(item.Points))
		author := cs.AuthorColor.Sprintf(item.User)
		timeAgo := cs.TimeAgoColor.Sprintf(item.TimeAgo)
		comments := cs.CommentsColor.Sprintf(strconv.Itoa(item.CommentsCount))

		firstLine := fmt.Sprintf("%s\t%s (%s)", rank, title, url)
		secondLine := fmt.Sprintf("\t%s points by %s %s| %s comments", points, author, timeAgo, comments)

		lines[index] = firstLine + "\n" + secondLine
	}
	return strings.Join(lines, "\n")
}
