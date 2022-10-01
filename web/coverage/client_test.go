package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

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
			w.WriteHeader(http.StatusOK)
		}),
	)

	ts.Close()

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

func TestStatusInternalServerError(t *testing.T) {

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

	result, _ := json.Marshal(
		map[string]string{
			"error": "something went wrong",
		},
	)

	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write(result)
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

	_, err := Client.FindUsers(Request)

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
			_, _ = w.Write(result)
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

func TestSearchServerSortAscID(t *testing.T) {

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
		Query:      "",
		OrderField: "Id",
		OrderBy:    1,
	}

	Response := SearchResponse{
		Users: []User{
			{
				ID:     0,
				Name:   "BoydWolf",
				Age:    22,
				About:  "Nulla cillum enim voluptate consequat laborum esse excepteur occaecat commodo nostrud excepteur ut cupidatat. Occaecat minim incididunt ut proident ad sint nostrud ad laborum sint pariatur. Ut nulla commodo dolore officia. Consequat anim eiusmod amet commodo eiusmod deserunt culpa. Ea sit dolore nostrud cillum proident nisi mollit est Lorem pariatur. Lorem aute officia deserunt dolor nisi aliqua consequat nulla nostrud ipsum irure id deserunt dolore. Minim reprehenderit nulla exercitation labore ipsum.\n",
				Gender: "male",
			},
		},
		NextPage: true,
	}

	result, err := Client.FindUsers(Request)

	if !reflect.DeepEqual(result.Users[0], Response.Users[0]) {
		t.Errorf("[TestSearchServerSortAscID] Response is wrong %v", result)
	}

	if err != nil {
		t.Errorf("[TestSearchServerSortAscID] expected error, got %v", err.Error())
	}
}

func TestSearchServerSortDescID(t *testing.T) {

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
		Query:      "",
		OrderField: "Id",
		OrderBy:    -1,
	}

	Response := SearchResponse{
		Users: []User{
			{
				ID:     1,
				Name:   "HildaMayer",
				Age:    21,
				About:  "Sit commodo consectetur minim amet ex. Elit aute mollit fugiat labore sint ipsum dolor cupidatat qui reprehenderit. Eu nisi in exercitation culpa sint aliqua nulla nulla proident eu. Nisi reprehenderit anim cupidatat dolor incididunt laboris mollit magna commodo ex. Cupidatat sit id aliqua amet nisi et voluptate voluptate commodo ex eiusmod et nulla velit.\n",
				Gender: "female",
			},
		},
		NextPage: true,
	}

	result, err := Client.FindUsers(Request)

	if !reflect.DeepEqual(result.Users[0], Response.Users[0]) {
		t.Errorf("[TestSearchServerSortDescID] Response is wrong %v", result)
	}

	if err != nil {
		t.Errorf("[TestSearchServerSortDescID] expected error, got %v", err.Error())
	}
}

func TestSearchServerSortAscAge(t *testing.T) {

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
		Query:      "",
		OrderField: "Age",
		OrderBy:    1,
	}

	Response := SearchResponse{
		Users: []User{
			{
				ID:     1,
				Name:   "HildaMayer",
				Age:    21,
				About:  "Sit commodo consectetur minim amet ex. Elit aute mollit fugiat labore sint ipsum dolor cupidatat qui reprehenderit. Eu nisi in exercitation culpa sint aliqua nulla nulla proident eu. Nisi reprehenderit anim cupidatat dolor incididunt laboris mollit magna commodo ex. Cupidatat sit id aliqua amet nisi et voluptate voluptate commodo ex eiusmod et nulla velit.\n",
				Gender: "female",
			},
		},
		NextPage: true,
	}

	result, err := Client.FindUsers(Request)

	if !reflect.DeepEqual(result.Users, Response.Users) {
		t.Errorf("[TestSearchServerSortAscAge] Response is wrong %v", result)
	}

	if err != nil {
		t.Errorf("[TestSearchServerSortAscAge] expected error, got %v", err.Error())
	}
}

func TestSearchServerSortDescAge(t *testing.T) {

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
		Query:      "",
		OrderField: "Age",
		OrderBy:    -1,
	}

	Response := SearchResponse{
		Users: []User{
			{
				ID:     0,
				Name:   "BoydWolf",
				Age:    22,
				About:  "Nulla cillum enim voluptate consequat laborum esse excepteur occaecat commodo nostrud excepteur ut cupidatat. Occaecat minim incididunt ut proident ad sint nostrud ad laborum sint pariatur. Ut nulla commodo dolore officia. Consequat anim eiusmod amet commodo eiusmod deserunt culpa. Ea sit dolore nostrud cillum proident nisi mollit est Lorem pariatur. Lorem aute officia deserunt dolor nisi aliqua consequat nulla nostrud ipsum irure id deserunt dolore. Minim reprehenderit nulla exercitation labore ipsum.\n",
				Gender: "male",
			},
		},
		NextPage: true,
	}

	result, err := Client.FindUsers(Request)

	if !reflect.DeepEqual(result.Users, Response.Users) {
		t.Errorf("[TestSearchServerSortDescAge] Response is wrong %v", result)
	}

	if err != nil {
		t.Errorf("[TestSearchServerSortDescAge] expected error, got %v", err.Error())
	}
}

func TestSearchServerSortAscName(t *testing.T) {

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
		Query:      "",
		OrderField: "Name",
		OrderBy:    1,
	}

	Response := SearchResponse{
		Users: []User{
			{
				ID:     0,
				Name:   "BoydWolf",
				Age:    22,
				About:  "Nulla cillum enim voluptate consequat laborum esse excepteur occaecat commodo nostrud excepteur ut cupidatat. Occaecat minim incididunt ut proident ad sint nostrud ad laborum sint pariatur. Ut nulla commodo dolore officia. Consequat anim eiusmod amet commodo eiusmod deserunt culpa. Ea sit dolore nostrud cillum proident nisi mollit est Lorem pariatur. Lorem aute officia deserunt dolor nisi aliqua consequat nulla nostrud ipsum irure id deserunt dolore. Minim reprehenderit nulla exercitation labore ipsum.\n",
				Gender: "male",
			},
		},
		NextPage: true,
	}

	result, err := Client.FindUsers(Request)

	if !reflect.DeepEqual(result.Users, Response.Users) {
		t.Errorf("[TestSearchServerSortAscName] Response is wrong %v", result)
	}

	if err != nil {
		t.Errorf("[TestSearchServerSortAscName] expected error, got %v", err.Error())
	}
}

func TestSearchServerSortDescName(t *testing.T) {

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
		Query:      "",
		OrderField: "Name",
		OrderBy:    -1,
	}

	Response := SearchResponse{
		Users: []User{
			{
				ID:     1,
				Name:   "HildaMayer",
				Age:    21,
				About:  "Sit commodo consectetur minim amet ex. Elit aute mollit fugiat labore sint ipsum dolor cupidatat qui reprehenderit. Eu nisi in exercitation culpa sint aliqua nulla nulla proident eu. Nisi reprehenderit anim cupidatat dolor incididunt laboris mollit magna commodo ex. Cupidatat sit id aliqua amet nisi et voluptate voluptate commodo ex eiusmod et nulla velit.\n",
				Gender: "female",
			},
		},
		NextPage: true,
	}

	result, err := Client.FindUsers(Request)

	if !reflect.DeepEqual(result.Users, Response.Users) {
		t.Errorf("[TestSearchServerSortDescName] Response is wrong %v", result)
	}

	if err != nil {
		t.Errorf("[TestSearchServerSortDescName] expected error, got %v", err.Error())
	}
}

func TestBadUnpackedJsonResponse(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
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
		Query:      "",
		OrderField: "Name",
		OrderBy:    -1,
	}

	_, err := Client.FindUsers(Request)

	returnError := fmt.Errorf("cant unpack result json: unexpected end of JSON input")

	if err.Error() != returnError.Error() {
		t.Errorf("[TestBadUnpackedJsonResponse] expected error, got %v", err.Error())
	}
}

func TestLenDataNotEqualLimit(t *testing.T) {

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
		Query:      "dommy",
		OrderField: "Name",
		OrderBy:    -1,
	}

	Response := SearchResponse{
		Users:    []User{},
		NextPage: false,
	}

	result, err := Client.FindUsers(Request)

	if !reflect.DeepEqual(result.Users, Response.Users) {
		t.Errorf("[TestLenDataNotEqualLimit] Response is wrong %v", result)
	}

	if err != nil {
		t.Errorf("[TestLenDataNotEqualLimit] expected error, got %v", err.Error())
	}
}

func TestBigOffset(t *testing.T) {

	ts := httptest.NewServer(
		http.HandlerFunc(SearchServer),
	)

	defer ts.Close()

	Client := SearchClient{
		URL:         ts.URL,
		AccessToken: "AccessToken",
	}
	Request := SearchRequest{
		Limit:      2,
		Offset:     34,
		Query:      "",
		OrderField: "Name",
		OrderBy:    -1,
	}

	Response := SearchResponse{
		Users: []User{
			{
				ID:     34,
				Name:   "KaneSharp",
				Age:    34,
				About:  "Lorem proident sint minim anim commodo cillum. Eiusmod velit culpa commodo anim consectetur consectetur sint sint labore. Mollit consequat consectetur magna nulla veniam commodo eu ut et. Ut adipisicing qui ex consectetur officia sint ut fugiat ex velit cupidatat fugiat nisi non. Dolor minim mollit aliquip veniam nostrud. Magna eu aliqua Lorem aliquip.\n",
				Gender: "male",
			},
		},
		NextPage: false,
	}

	result, err := Client.FindUsers(Request)

	if !reflect.DeepEqual(result.Users, Response.Users) {
		t.Errorf("[TestBigOffset] Response is wrong %v", result)
	}

	if err != nil {
		t.Errorf("[TestBigOffset] expected error, got %v", err.Error())
	}
}

func TestSearchServerBadXml(t *testing.T) {

	dataset = "bad_dataset.xml"

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
		Query:      "",
		OrderField: "Name",
		OrderBy:    1,
	}

	_, err := Client.FindUsers(Request)
	returnError := fmt.Errorf("cant unpack error json: json: cannot unmarshal array into Go value of type main.SearchErrorResponse")

	if err.Error() != returnError.Error() {
		t.Errorf("[TestSearchServerBadXml] expected error, got %v", err.Error())
	}
}
