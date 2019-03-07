package news

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const apiURL = "https://api.hnpwa.com/"

type Option func(*Client) error

type Client struct {
	baseUrl *url.URL
	client  *http.Client
}

func NewClient(opts ...Option) (*Client, error) {
	u, _ := url.Parse(apiURL)
	c := &Client{
		baseUrl: u,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	for _, option := range opts {
		err := option(c)
		if err != nil {
			return nil, err
		}
	}

	return c, nil
}

func BaseURL(baseUrl string) Option {
	return func(c *Client) error {
		u, err := url.Parse(baseUrl)
		if err != nil {
			return err
		}
		c.baseUrl = u
		return nil
	}
}

func HttpClient(client *http.Client) Option {
	return func(c *Client) error {
		c.client = client
		return nil
	}
}

func (c *Client) GetFeed(ctx context.Context, kind FeedKind, page int) (*Feed, error) {
	var items []*Item

	rel := &url.URL{Path: fmt.Sprintf("/v0/%s/%d.json", kind, page)}
	u := c.baseUrl.ResolveReference(rel)

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	_, err = c.runHttpRequest(ctx, req, &items)
	feed := &Feed{Kind: kind, Page: page, Items: items}
	return feed, err
}

func (c *Client) GetDetail(ctx context.Context, id int) (*Detail, error) {
	var detail Detail

	rel := &url.URL{Path: fmt.Sprintf("/v0/item/%d.json", id)}
	u := c.baseUrl.ResolveReference(rel)

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	_, err = c.runHttpRequest(ctx, req, &detail)
	return &detail, err
}

// Run the given http request and deserialize the response into the given interface
func (c *Client) runHttpRequest(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
	req = req.WithContext(ctx)
	resp, err := c.client.Do(req)
	if err != nil {
		select {
		case <-ctx.Done():
			return resp, ctx.Err()
		default:
			return resp, err
		}
	}
	defer resp.Body.Close()

	if err := checkResponse(resp); err != nil {
		return resp, err
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp, err
	}

	err = json.Unmarshal(bytes, &v)
	return resp, err
}

func checkResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	return errors.New(fmt.Sprintf("HTTP %d: %s", r.StatusCode, string(data)))
}
