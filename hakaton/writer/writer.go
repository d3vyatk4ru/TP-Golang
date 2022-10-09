package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

type RecordHandler struct {
	file os.File
}

func (record *RecordHandler) handler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		fmt.Println("no POST")
	}

	dataByte, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		fmt.Println("bad read handle: ", err)
	}

	data := string(dataByte[:])
	data = data[:len(data)-1]
	now := time.Now()
	sec := now.Unix()
	strTime := strconv.FormatInt(sec, 10)
	data = data + ",\"Time\":\"" + strTime + "\"}" + "\n"
	fmt.Println(data)
	_, err = record.file.WriteString(data)
	if err != nil {
		fmt.Println("Errr", err)
	}
}

func main() {

	port := flag.String("p", "6001", "a string")
	fileName := flag.String("f", "writer1.txt", "a string")
	flag.Parse()

	f, err := os.Create(*fileName)
	if err != nil {
		fmt.Println("bad create file", err)
	}

	record := RecordHandler{
		file: *f,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", record.handler)

	addr := ":" + *port
	server := http.Server{
		Addr:    addr,
		Handler: mux,
	}

	err = server.ListenAndServe()
}
