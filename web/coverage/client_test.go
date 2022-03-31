package main

import (
	"errors"
	"fmt"
	"testing"
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
				URL:         "www.mail.ru",
				AccessToken: "xyz",
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

		if errors.Is(err, item.returnValue) {
			t.Errorf("[TestFindUsersBadLimit](%d) expected error, got %e", num, err)
		}

	}

}

func TestFindUsersBadOffsetAndDecreaseLimit(t *testing.T) {

	cases := []TestCase{
		{
			Client: SearchClient{
				URL:         "www.mail.ru",
				AccessToken: "xyz",
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

		if errors.Is(err, item.returnValue) {
			t.Errorf("[TestFindUsersBadOffsetAnd...](%d) expected error, got %e", num, err)
		}
	}
}

func TestFindUsersBadDo(t *testing.T) {

	cases := []TestCase{
		{
			Client: SearchClient{
				URL:         "www.mail.ru",
				AccessToken: "xyz",
			},
			Request: SearchRequest{
				Limit:      1,
				Offset:     0,
				Query:      "abc",
				OrderField: "cba",
				OrderBy:    0,
			},
			returnValue: fmt.Errorf("GET www.mail.ru?limit=2&offset=0&order_by=0&order_field=cba&query=abc")},
	}

	for num, item := range cases {

		_, err := item.Client.FindUsers(item.Request)

		if errors.Is(err, item.returnValue) {
			t.Errorf("[TestFindUsersBadDo](%d) expected error, got %e", num, err)
		}
	}
}

func TestFindUsersBadDoTimeout(t *testing.T) {

}

func FindUsers(searchRequest *SearchRequest) {
	panic("unimplemented")
}

// тут писать код тестов
