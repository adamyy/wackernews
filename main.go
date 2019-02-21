package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/adamyy/hackernews/api"
	"github.com/adamyy/hackernews/feed"
	"github.com/adamyy/hackernews/view"
)

func main() {
	f, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}
	p, err := strconv.Atoi(os.Args[2])
	if err != nil {
		panic(err)
	}

	ft := feed.Type(f)
	items, err := api.GetFeeds(ft, p)
	if err != nil {
		panic(err)
	}
	fv := view.NewFeedView(ft, items, p)
	fmt.Println(fv.Render())
}
