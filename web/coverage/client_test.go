package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

type TestCase struct {
	Request     SearchRequest
	Response    SearchResponse
	returnValue error
	Client      SearchClient
}

func TestFindUsersBadLimit(t *testing.T) {

	Client := SearchClient{
		URL:         "localhost:8080",
		AccessToken: "AccessToken",
	}
	Request := SearchRequest{
		Limit:      -1,
		Offset:     1,
		Query:      "abc",
		OrderField: "Name",
		OrderBy:    0,
	}

	returnError := fmt.Errorf("limit must be > 0")

	_, err := Client.FindUsers(Request)

	if err.Error() != returnError.Error() {
		t.Errorf("[TestFindUsersBadLimit] expected error, got %v", err.Error())
	}
}

func TestFindUsersBadOffsetAndDecreaseLimit(t *testing.T) {

	Client := SearchClient{
		URL:         "localhost:8080",
		AccessToken: "AccessToken",
	}

	Request := SearchRequest{
		Limit:      26,
		Offset:     -1,
		Query:      "abc",
		OrderField: "Name",
		OrderBy:    0,
	}

	returnError := fmt.Errorf("offset must be > 0")

	_, err := Client.FindUsers(Request)

	if err.Error() != returnError.Error() {
		t.Errorf("[TestFindUsersBadOffsetAnd...] expected error, got %v", err.Error())
	}
}

func TestFindUsersBadDo(t *testing.T) {

	Client := SearchClient{
		URL:         "www.mail.ru",
		AccessToken: "AccessToken",
	}
	Request := SearchRequest{
		Limit:      1,
		Offset:     0,
		Query:      "abc",
		OrderField: "Name",
		OrderBy:    0,
	}
	returnError := fmt.Errorf(`unknown error Get "www.mail.ru?limit=2&offset=0&order_by=0&order_field=Name&query=abc": unsupported protocol scheme ""`)

	_, err := Client.FindUsers(Request)

	if err.Error() != returnError.Error() {
		t.Errorf("[TestFindUsersBadDo] %v", err.Error())
	}

}

func TestFindUsersBadDoTimeout(t *testing.T) {

	// создаем липовый сервер
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(time.Second)
		}),
	)

	defer ts.Close()

	Client := SearchClient{
		URL:         ts.URL,
		AccessToken: "AccessToken",
	}
	Request := SearchRequest{
		Limit:      1,
		Offset:     0,
		Query:      "abc",
		OrderField: "Name",
		OrderBy:    0,
	}
	returnError := fmt.Errorf("timeout for limit=2&offset=0&order_by=0&order_field=Name&query=abc")

	_, err := Client.FindUsers(Request)

	if err.Error() != returnError.Error() {
		t.Errorf("[TestFindUsersBadDoTimeout] expected error, got %v", err.Error())
	}
}

func TestFindUsersStatusUnauthorized(t *testing.T) {

	ts := httptest.NewServer(
		http.HandlerFunc(SearchServer),
	)

	defer ts.Close()

	Client := SearchClient{
		URL:         ts.URL,
		AccessToken: "BadAccessToken",
	}
	Request := SearchRequest{
		Limit:      1,
		Offset:     0,
		Query:      "abc",
		OrderField: "Name",
		OrderBy:    0,
	}
	returnError := fmt.Errorf("bad AccessToken")

	_, err := Client.FindUsers(Request)

	if err.Error() != returnError.Error() {
		t.Errorf("[TestFindUsersStatusUnauthorized] expected error, got %v", err.Error())
	}
}

func StatusInternalServerError(t *testing.T) {

	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			}),
	)

	defer ts.Close()

	Client := SearchClient{
		URL:         ts.URL,
		AccessToken: "AccessToken",
	}
	Request := SearchRequest{
		Limit:      1,
		Offset:     0,
		Query:      "abc",
		OrderField: "Name",
		OrderBy:    0,
	}
	returnError := fmt.Errorf("SearchServer fatal error")

	_, err := Client.FindUsers(Request)

	if err.Error() != returnError.Error() {
		t.Errorf("[TestFindUsers] expected error, got %e", err)
	}
}

func TestStatusBadRequestBadUnpackedJson(t *testing.T) {

	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
		}),
	)

	defer ts.Close()

	Client := SearchClient{
		URL:         ts.URL,
		AccessToken: "AccessToken",
	}
	Request := SearchRequest{
		Limit:      1,
		Offset:     0,
		Query:      "abc",
		OrderField: "Name",
		OrderBy:    0,
	}
	returnError := fmt.Errorf("cant unpack error json: unexpected end of JSON input")

	_, err := Client.FindUsers(Request)

	if err.Error() != returnError.Error() {
		t.Errorf("[TestStatusBadRequestBadUnpackedJson] expected error, got %v", err.Error())
	}
}

func TestStatusBadRequest(t *testing.T) {

	result, err := json.Marshal(
		map[string]string{
			"error": "something went wrong",
		},
	)

	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(result)
		}),
	)

	defer ts.Close()

	Client := SearchClient{
		URL:         ts.URL,
		AccessToken: "AccessToken",
	}
	Request := SearchRequest{
		Limit:      1,
		Offset:     0,
		Query:      "abc",
		OrderField: "Name",
		OrderBy:    0,
	}

	returnError := fmt.Errorf("unknown bad request error: something went wrong")

	_, err = Client.FindUsers(Request)

	if err.Error() != returnError.Error() {
		t.Errorf("[TestStatusBadRequestBadUnpackedJson] expected error, got %v", err.Error())
	}
}
func TestStatusBadRequestErrorBadOrderField(t *testing.T) {

	result, _ := json.Marshal(
		SearchErrorResponse{Error: `OrderField invalid`},
	)

	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(result)
		}),
	)

	defer ts.Close()

	Client := SearchClient{
		AccessToken: "AccessToken",
		URL:         ts.URL,
	}
	Request := SearchRequest{
		Limit:      1,
		Offset:     0,
		Query:      "abc",
		OrderField: "InvalidOrderField",
		OrderBy:    1,
	}
	returnError := fmt.Errorf("OrderField InvalidOrderField invalid")

	_, err := Client.FindUsers(Request)

	if err.Error() != returnError.Error() {
		t.Errorf("[TestStatusBadRequestErrorBadOrderField] expected error, got %v", err.Error())
	}
}

// func TestBadUnpackedJson(t *testing.T) {

// 	result, _ := json.Marshal([]User{})

// 	ts := httptest.NewServer(
// 		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			w.Write(result)
// 		}),
// 	)

// 	defer ts.Close()

// 	Client := SearchClient{
// 		URL:         ts.URL,
// 		AccessToken: "AccessToken",
// 	}
// 	Request := SearchRequest{
// 		Limit:      1,
// 		Offset:     0,
// 		Query:      "abc",
// 		OrderField: "Name",
// 		OrderBy:    0,
// 	}

// 	returnValue := fmt.Errorf("cant unpack result json: unexpected end of JSON input")

// 	_, err := Client.FindUsers(Request)

// 	if err.Error() != returnValue.Error() {
// 		t.Errorf("[TestBadUnpackedJson] expected error, got %v", err.Error())
// 	}
// }

func TestFindUser(t *testing.T) {

	ts := httptest.NewServer(
		http.HandlerFunc(SearchServer),
	)

	defer ts.Close()

	Client := SearchClient{
		URL:         ts.URL,
		AccessToken: "AccessToken",
	}
	Request := SearchRequest{
		Limit:      1,
		Offset:     0,
		Query:      "abc",
		OrderField: "Name",
		OrderBy:    0,
	}
	Response := SearchResponse{
		Users: []User{
			{
				ID:     999,
				Name:   "Default",
				Age:    0,
				About:  "",
				Gender: "male",
			},
		},
		NextPage: false,
	}

	result, err := Client.FindUsers(Request)

	if !reflect.DeepEqual(result.Users, Response.Users) {
		t.Errorf("[TestBadUnpackedJson] Response is wrong %v", result)
	}

	if err != nil {
		t.Errorf("[TestBadUnpackedJson] expected error, got %v", err.Error())
	}
}
