package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/adamyy/hackernews/api"
	"github.com/adamyy/hackernews/feed"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		cmds := strings.Fields(text)
		switch cmds[0] {
		case "feed":
			t := feed.TypeOf(cmds[1])
			p, err := strconv.Atoi(cmds[2])
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			feeds, err := api.GetFeeds(t, p)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			for _, feed := range feeds {
				fmt.Printf("%+v\n", feed)
			}
		case "detail":
			id, err := strconv.Atoi(cmds[1])
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			detail, err := api.GetDetail(id)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			fmt.Printf("%+v\n", detail)
		case "exit":
			break
		}
	}
}
