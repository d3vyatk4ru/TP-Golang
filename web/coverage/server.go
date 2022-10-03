package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

var dataset = "dataset.xml"

type Root struct {
	XMLName xml.Name `xml:"root"`
	Rows    []Row    `xml:"row"`
}

type Row struct {
	XMLName   xml.Name `xml:"row"`
	ID        int      `xml:"id"`
	Age       int      `xml:"age"`
	FirstName string   `xml:"first_name"`
	LastName  string   `xml:"last_name"`
	Gender    string   `xml:"gender"`
	About     string   `xml:"about"`
}

//nolint:gocyclo
func SearchServer(w http.ResponseWriter, r *http.Request) {

	var record Root
	var data []Row
	var user = make([]User, 0)

	accessToken := r.Header.Get("AccessToken")

	if accessToken != "AccessToken" {
		fmt.Println("Bad AccessToken")
		w.WriteHeader(http.StatusUnauthorized)
		emptyUser, err := json.Marshal([]User{})
		if err != nil {
			fmt.Println("Bad pack json")
		}
		_, err = w.Write(emptyUser)
		if err != nil {
			fmt.Println("Bad write response")
		}
		return
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))

	if err != nil {
		fmt.Println("Bad convert to int")
		return
	}

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		fmt.Println("Bad convert to int")
		return
	}

	OrderBy, err := strconv.Atoi(r.URL.Query().Get("order_by"))
	if err != nil {
		fmt.Println("Bad convert to int")
		return
	}

	OrderField := r.URL.Query().Get("order_field")

	query := r.URL.Query().Get("query")

	xmlDsc, err := ioutil.ReadFile(dataset)

	if err != nil {
		fmt.Println("Trouble with xml file: ", err.Error())
		fmt.Println("Inside xml error")
		w.WriteHeader(http.StatusBadRequest)
		data, err2 := json.Marshal([]User{})
		if err2 != nil {
			fmt.Println("Bad pack json")
		}
		_, err2 = w.Write(data)
		if err2 != nil {
			fmt.Println("Bad write response")
		}
		return
	}

	err = xml.Unmarshal(xmlDsc, &record)

	if err != nil {
		fmt.Println("cant unpack result xml: ", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		data, err2 := json.Marshal([]User{})
		if err2 != nil {
			fmt.Println("Bad pack json")
		}
		_, err2 = w.Write(data)
		if err2 != nil {
			fmt.Println("Bad write response")
		}
		return
	}

	if offset+limit > len(record.Rows) {
		data = record.Rows[offset:]
	} else {
		data = record.Rows[offset : offset+limit]
	}

	switch OrderField {
	case "Id":

		if OrderBy == -1 {
			sort.Slice(data, func(i, j int) bool {
				return data[i].ID > data[j].ID
			})
		}

		if OrderBy == 1 {
			sort.Slice(data, func(i, j int) bool {
				return data[i].ID < data[j].ID
			})
		}
	case "Age":

		if OrderBy == -1 {
			sort.Slice(data, func(i, j int) bool {
				return data[i].Age > data[j].Age
			})
		}

		if OrderBy == 1 {
			sort.Slice(data, func(i, j int) bool {
				return data[i].Age < data[j].Age
			})
		}
	default:

		if OrderBy == -1 {
			sort.Slice(data, func(i, j int) bool {
				return data[i].FirstName+data[i].LastName >
					data[j].FirstName+data[j].LastName
			})
		}

		if OrderBy == 1 {
			sort.Slice(data, func(i, j int) bool {
				return data[i].FirstName+data[i].LastName <
					data[j].FirstName+data[j].LastName
			})
		}
	}

	for _, row := range data {

		if strings.Contains(row.About, query) ||
			strings.Contains(row.FirstName+row.LastName, query) ||
			strings.Contains(row.Gender, query) {
			user = append(user,
				User{
					ID:     row.ID,
					Name:   row.FirstName + row.LastName,
					Age:    row.Age,
					About:  row.About,
					Gender: row.Gender,
				},
			)
		}
	}

	result, err := json.Marshal(user)

	if err != nil {
		fmt.Println("cant pack result json: ", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		data, err2 := json.Marshal([]User{})
		if err2 != nil {
			fmt.Println("Bad pack json")
		}
		_, err2 = w.Write(data)
		if err2 != nil {
			fmt.Println("Bad write response")
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(result)

	if err != nil {
		fmt.Println("Bad write response")
	}
}
