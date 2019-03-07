package news

type Detail struct {
	Id            int64     `json:"id"`
	Title         string    `json:"title"`
	Points        int       `json:"points"`
	User          string    `json:"user"`
	Time          int       `json:"time"`
	TimeAgo       string    `json:"time_ago"`
	Content       string    `json:"content"`
	Deleted       bool      `json:"deleted"`
	Dead          bool      `json:"dead"`
	Type          string    `json:"type"`
	Url           string    `json:"url"`
	Domain        string    `json:"domain"`
	Comments      []*Detail `json:"comments"`
	Level         int       `json:"level"`
	CommentsCount int       `json:"comments_count"`
}
