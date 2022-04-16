package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type TestCase struct {
	Request     SearchRequest
	returnValue error
	Client      SearchClient
}

func TestFindUsersBadLimit(t *testing.T) {

	cases := []TestCase{
		{
			Client: SearchClient{
				URL:         "localhost:8080",
				AccessToken: "AccessToken",
			},
			Request: SearchRequest{
				Limit:      -1,
				Offset:     1,
				Query:      "abc",
				OrderField: "cba",
				OrderBy:    0,
			},
			returnValue: fmt.Errorf("limit must be > 0")},
	}

	for num, item := range cases {

		_, err := item.Client.FindUsers(item.Request)

		if err.Error() != item.returnValue.Error() {
			t.Errorf("[TestFindUsersBadLimit](%d) expected error, got %e", num, err)
		}
	}
}

func TestFindUsersBadOffsetAndDecreaseLimit(t *testing.T) {

	cases := []TestCase{
		{
			Client: SearchClient{
				URL:         "localhost:8080",
				AccessToken: "AccessToken",
			},
			Request: SearchRequest{
				Limit:      26,
				Offset:     -1,
				Query:      "abc",
				OrderField: "cba",
				OrderBy:    0,
			},
			returnValue: fmt.Errorf("offset must be > 0")},
	}

	for num, item := range cases {

		_, err := item.Client.FindUsers(item.Request)

		if err.Error() != item.returnValue.Error() {
			t.Errorf("[TestFindUsersBadOffsetAnd...](%d) expected error, got %e", num, err)
		}
	}
}

func TestFindUsersBadDo(t *testing.T) {

	cases := []TestCase{
		{
			Client: SearchClient{
				URL:         "www.mail.ru",
				AccessToken: "AccessToken",
			},
			Request: SearchRequest{
				Limit:      1,
				Offset:     0,
				Query:      "abc",
				OrderField: "cba",
				OrderBy:    0,
			},
			returnValue: fmt.Errorf("unknown error Get www.mail.ru?limit=2&offset=0&order_by=0&order_field=cba&query=abc: unsupported protocol scheme ")},
	}

	for num, item := range cases {

		_, err := item.Client.FindUsers(item.Request)

		if err.Error() != item.returnValue.Error() {
			t.Errorf("[TestFindUsersBadDo](%d) expected error, got %e", num, err)
		}
	}
}

func TestFindUsersBadDoTimeout(t *testing.T) {

	// создаем липовый сервер
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Second)
	}))

	cases := []TestCase{
		{
			Client: SearchClient{
				URL:         ts.URL,
				AccessToken: "AccessToken",
			},
			Request: SearchRequest{
				Limit:      1,
				Offset:     0,
				Query:      "abc",
				OrderField: "cba",
				OrderBy:    0,
			},
			returnValue: fmt.Errorf("timeout for limit=2&offset=0&order_by=0&order_field=cba&query=abc")},
	}

	for num, item := range cases {

		_, err := item.Client.FindUsers(item.Request)

		if err.Error() != item.returnValue.Error() {
			t.Errorf("[TestFindUsersBadDoTimeout](%d) expected error, got %e", num, err)
		}
	}
}

func TestFindUsersStatusUnauthorized(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(SearchServer))

	cases := []TestCase{
		{
			Client: SearchClient{
				URL:         ts.URL,
				AccessToken: "BadAccessToken",
			},
			Request: SearchRequest{
				Limit:      1,
				Offset:     0,
				Query:      "abc",
				OrderField: "cba",
				OrderBy:    0,
			},
			returnValue: fmt.Errorf("bad AccessToken")},
	}

	for num, item := range cases {

		_, err := item.Client.FindUsers(item.Request)

		if err.Error() != item.returnValue.Error() {
			t.Errorf("[TestFindUsersStatusUnauthorized](%d) expected error, got %e", num, err)
		}
	}
}

func TestFindUsers(t *testing.T) {

	ts := httptest.NewServer((http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})))

	cases := []TestCase{
		{
			Client: SearchClient{
				URL:         ts.URL,
				AccessToken: "AccessToken",
			},
			Request: SearchRequest{
				Limit:      1,
				Offset:     0,
				Query:      "abc",
				OrderField: "cba",
				OrderBy:    0,
			},
			returnValue: fmt.Errorf("SearchServer fatal error")},
	}

	httptest.NewServer(http.HandlerFunc(SearchServer))

	for _, item := range cases {

		_, err := item.Client.FindUsers(item.Request)

		if err.Error() != item.returnValue.Error() {

		}
	}
}
