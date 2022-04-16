package main

import (
	"net/http"
)

func SearchServer(w http.ResponseWriter, r *http.Request) {

	accessToken := r.Header.Get("AccessToken")

	if accessToken != "AccessToken" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// query := r.FormValue("query")
}
