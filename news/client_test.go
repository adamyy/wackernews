package news

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	mux    *http.ServeMux
	server *httptest.Server
	client *Client
)

func setup() func() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	client, _ = NewClient(BaseURL(server.URL))

	return func() { server.Close() }
}

func loadFixture(name string) string {
	b, err := ioutil.ReadFile("testdata/fixtures/" + name)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func TestApi_GetFeed(t *testing.T) {
	teardown := setup()
	defer teardown()

	tests := []struct {
		kind       FeedKind
		page       int
		err        error
		fixture    string
		statusCode int
	}{
		{
			kind:       KindNews,
			page:       1,
			fixture:    "get_feed_news_ok",
			statusCode: http.StatusOK,
		},
		{
			kind:       KindAsk,
			page:       0,
			err:        errors.New("HTTP 500: Error: could not handle the request"),
			statusCode: http.StatusInternalServerError,
			fixture:    "get_feed_news_server_error",
		},
	}

	for _, test := range tests {
		mux.HandleFunc(fmt.Sprintf("/%s/%d.json", test.kind, test.page),
			func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "GET")
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(test.statusCode)
				fmt.Fprintf(w, loadFixture(test.fixture))
			})

		feed, err := client.GetFeed(context.Background(), test.kind, test.page)
		if err != nil {
			if test.err != nil {
				assert.Equal(t, test.err, err)
				continue
			}
			t.Error(err)
		} else {
			assert.Equal(t, test.kind, feed.Kind)
			assert.Equal(t, test.page, feed.Page)
			assert.Equal(t, 30, len(feed.Items))
		}
	}
}

func TestApi_GetDetail(t *testing.T) {
	teardown := setup()
	defer teardown()

	tests := []struct {
		id           int
		err          error
		fixture      string
		statusCode   int
		errorMessage string
	}{
		{
			id:         13831370,
			fixture:    "get_detail_ok",
			statusCode: http.StatusOK,
		},
		{
			id:         0,
			statusCode: http.StatusInternalServerError,
			fixture:    "get_detail_server_error",
			err:        errors.New("HTTP 500: Error: could not handle the request"),
		},
	}

	for _, test := range tests {
		mux.HandleFunc(fmt.Sprintf("/item/%d.json", test.id),
			func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "GET")
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(test.statusCode)
				fmt.Fprintf(w, loadFixture(test.fixture))
			})

		detail, err := client.GetDetail(context.Background(), test.id)
		if err != nil {
			if test.err != nil {
				assert.Equal(t, test.err, err)
				continue
			}
			t.Error(err)
		} else {
			assert.Equal(t, test.id, detail.Id)
		}
	}
}

func testMethod(t *testing.T, r *http.Request, want string) {
	if got := r.Method; got != want {
		t.Errorf("Request method: %v, want %v", got, want)
	}
}
