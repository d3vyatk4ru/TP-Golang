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
	XMLName    xml.Name `xml:"row"`
	Id         int      `xml:"id"`
	Age        int      `xml:"age"`
	First_name string   `xml:"first_name"`
	Last_name  string   `xml:"last_name"`
	Gender     string   `xml:"gender"`
	About      string   `xml:"about"`
}

func SearchServer(w http.ResponseWriter, r *http.Request) {

	var record Root
	var data []Row
	var user = make([]User, 0)

	accessToken := r.Header.Get("AccessToken")

	if accessToken != "AccessToken" {
		fmt.Println("Bad AccessToken")
		w.WriteHeader(http.StatusUnauthorized)
		emptyUser, _ := json.Marshal([]User{})
		_, _ = w.Write(emptyUser)
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	order_by, _ := strconv.Atoi(r.URL.Query().Get("order_by"))

	order_field := r.URL.Query().Get("order_field")

	query := r.URL.Query().Get("query")

	xmlDsc, err := ioutil.ReadFile(dataset)

	if err != nil {
		fmt.Println("Trouble with xml file: ", err.Error())
		fmt.Println("Inside xml error")
		w.WriteHeader(http.StatusBadRequest)
		data, _ := json.Marshal([]User{})
		_, _ = w.Write(data)
		return
	}

	err = xml.Unmarshal(xmlDsc, &record)

	if err != nil {
		fmt.Println("cant unpack result xml: ", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		data, _ := json.Marshal([]User{})
		_, _ = w.Write([]byte(data))
		return
	}

	if offset+limit > len(record.Rows) {
		data = record.Rows[offset:]
	} else {
		data = record.Rows[offset : offset+limit]
	}

	if order_field == "Id" {

		if order_by == -1 {
			sort.Slice(data, func(i, j int) bool {
				return data[i].Id > data[j].Id
			})
		}

		if order_by == 1 {
			sort.Slice(data, func(i, j int) bool {
				return data[i].Id < data[j].Id
			})
		}
	} else if order_field == "Age" {

		if order_by == -1 {
			sort.Slice(data, func(i, j int) bool {
				return data[i].Age > data[j].Age
			})
		}

		if order_by == 1 {
			sort.Slice(data, func(i, j int) bool {
				return data[i].Age < data[j].Age
			})
		}
	} else {

		if order_by == -1 {
			sort.Slice(data, func(i, j int) bool {
				return data[i].First_name+data[i].Last_name >
					data[j].First_name+data[j].Last_name
			})
		}

		if order_by == 1 {
			sort.Slice(data, func(i, j int) bool {
				return data[i].First_name+data[i].Last_name <
					data[j].First_name+data[j].Last_name
			})
		}
	}

	for _, row := range data {

		if strings.Contains(row.About, query) ||
			strings.Contains(row.First_name+row.Last_name, query) ||
			strings.Contains(row.Gender, query) {
			user = append(user,
				User{
					ID:     row.Id,
					Name:   row.First_name + row.Last_name,
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
		data, _ := json.Marshal([]User{})
		_, _ = w.Write(data)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(result))
}
