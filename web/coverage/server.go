package main

import (
	"encoding/json"
	"net/http"
)

func SearchServer(w http.ResponseWriter, r *http.Request) {

	accessToken := r.Header.Get("AccessToken")

	if accessToken != "AccessToken" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// query := r.FormValue("query")

	result, _ := json.Marshal(
		[]User{
			{
				ID:     999,
				Name:   "Default",
				Age:    0,
				About:  "",
				Gender: "male",
			},
		},
	)

	w.Write(result)

	return
}
