package news

import "strings"

type FeedKind int

type Feed struct {
	Kind  FeedKind
	Page  int
	Items []*Item
}

type Item struct {
	Id            int    `json:"id"`
	Title         string `json:"title"`
	Points        int    `json:"points"`
	User          string `json:"user"`
	Time          int    `json:"time"`
	TimeAgo       string `json:"timeago"`
	CommentsCount int    `json:"comments_count"`
	Type          string `json:"type"`
	Url           string `json:"url"`
	Domain        string `json:"domain"`
}

const (
	KindNews FeedKind = iota
	KindNewest
	KindAsk
	KindShow
	KindJobs
)

var (
	AllFeedKinds = [...]FeedKind{KindNews, KindNewest, KindAsk, KindShow, KindJobs}

	feedText = map[FeedKind]string{
		KindNews:   "news",
		KindNewest: "newest",
		KindAsk:    "ask",
		KindShow:   "show",
		KindJobs:   "jobs",
	}

	stringToKind = map[string]FeedKind{
		"news":   KindNews,
		"newest": KindNewest,
		"ask":    KindAsk,
		"show":   KindShow,
		"jobs":   KindJobs,
	}

	feedReadableText = map[FeedKind]string{
		KindNews:   "News",
		KindNewest: "Newest",
		KindAsk:    "Ask",
		KindShow:   "Show",
		KindJobs:   "Jobs",
	}
)

func (k FeedKind) String() string {
	return feedText[k]
}

func (k FeedKind) ReadableString() string {
	return feedReadableText[k]
}

func ToFeedKind(s string) FeedKind {
	return stringToKind[strings.ToLower(s)]
}
