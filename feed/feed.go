package feed

import "strings"

type Type int

var StringToType = map[string]Type{
	"news":   TypeNews,
	"newest": TypeNewest,
	"ask":    TypeAsk,
	"show":   TypeShow,
	"jobs":   TypeJobs,
}

func TypeOf(s string) Type {
	return StringToType[strings.ToLower(s)]
}

const (
	TypeNews Type = iota
	TypeNewest
	TypeAsk
	TypeShow
	TypeJobs
)

func (f Type) String() string {
	return [...]string{"news", "newest", "ask", "show", "jobs"}[f]
}

func (f Type) ReadableString() string {
	return [...]string{"News", "Newest", "Ask", "Show", "Jobs"}[f]
}

var AllTypes = [...]Type{TypeNews, TypeNewest, TypeAsk, TypeShow, TypeJobs}

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

type Detail struct {
	Id          int    `json:"id"`          // The item's unique id.
	Deleted     bool   `json:"deleted"`     // true if the item is deleted.
	Type        string `json:"type"`        // The type of item. One of "job", "story", "comment", "poll", or "pollopt".
	By          string `json:"by"`          // The username of the item's author.
	Time        int    `json:"time"`        // Creation date of the item, in Unix Time.
	Text        string `json:"text"`        // The comment, story or poll text. HTML.
	Dead        bool   `json:"dead"`        // true if the item is dead.
	Parent      int    `json:"parent"`      // The comment's parent: either another comment or the relevant story.
	Poll        int    `json:"poll"`        // The pollopt's associated poll.
	Kids        []int  `json:"kids"`        // The ids of the item's comments, in ranked display order.
	Url         string `json:"url"`         // The URL of the story.
	Score       int    `json:"score"`       // The story's score, or the votes for a pollopt.
	Title       string `json:"title"`       // The title of the story, poll or job.
	Parts       []int  `json:"parts"`       // A list of related pollopts, in display order.
	Descendants int    `json:"descendants"` // In the case of stories or polls, the total comment count.
}
