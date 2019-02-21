package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/adamyy/hackernews/feed"
)

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

func GetFeeds(feedType feed.Type, page int) ([]feed.Item, error) {
	var feeds []feed.Item

	req, err := feedsRequest(feedType, page)
	if err != nil {
		return feeds, err
	}

	err = runHttpRequest(req, &feeds)

	return feeds, err
}

func feedsRequest(feedType feed.Type, page int) (*http.Request, error) {
	return http.NewRequest(http.MethodGet, fmt.Sprintf("https://api.hnpwa.com/v0/%s/%d.json", feedType, page), nil)
}

func GetDetail(id int) (*feed.Detail, error) {
	detail := new(feed.Detail)

	req, err := detailRequest(id)
	if err != nil {
		return detail, nil
	}

	err = runHttpRequest(req, &detail)
	return detail, nil
}

func detailRequest(id int) (*http.Request, error) {
	return http.NewRequest(http.MethodGet, fmt.Sprintf("https://api.hnpwa.com/v0/item/%d.json", id), nil)
}

// Run the given http request and deserialize the response into the given interface
func runHttpRequest(r *http.Request, v interface{}) error {
	resp, err := httpClient.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	json.Unmarshal(bytes, &v)
	return nil
}
