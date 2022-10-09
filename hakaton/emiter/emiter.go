package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"encoding/json"
	"flag"
	"fmt"
	"hash/crc32"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var writerUrl string = "http://127.0.0.1"

type Record struct {
	ID        uint32
	ASIN      string
	Title     string
	Group     string
	Salesrank string
}

func splitPart(c rune) bool {
	return c == ' '
}

func ParseData(s string) (Record, error) {
	ParsedData := Record{}
	SplittedData := strings.Split(s, "\n")
	for _, row := range SplittedData {
		parts := strings.FieldsFunc(strings.TrimSpace(row), splitPart)
		switch parts[0] {
		case "Id:":
			id, err := strconv.ParseUint(string(parts[1]), 10, 64)
			if err != nil {
				fmt.Println("cannot read id")
			}
			ParsedData.ID = uint32(id)
		case "ASIN:":
			ParsedData.ASIN = parts[1]
		case "title:":
			ParsedData.Title = strings.Join(parts[1:], " ")
		case "group:":
			ParsedData.Group = parts[1]
		case "salesrank:":
			ParsedData.Salesrank = parts[1]
		case "similar:":
			break
		case "discontinued":
			return ParsedData, nil
		default:
			return ParsedData, fmt.Errorf("Parsing error")
		}
	}

	return ParsedData, nil
}

func fileParser(filename string, clients []http.Client, ports []string) error {

	file, err := os.Open(filename)
	defer file.Close()

	if err != nil {
		log.Fatal(err)
		return err
	}

	archive, err := gzip.NewReader(file)

	if err != nil {
		log.Fatal(err)
		return err
	}

	scanner := bufio.NewScanner(archive)

	for scanner.Scan() {
		split := strings.Split(string(scanner.Text()), ":")

		if string(split[0]) == "Id" {
			data := split[0] + ": " + split[1] + "\n"
			for scanner.Scan() {
				if string(scanner.Text()) != "" {
					data += string(scanner.Text()) + "\n"
				} else {
					break
				}
			}

			record, _ := ParseData(data)
			writerId := getReaderIDforRecord(&record, uint32(len(clients)))

			postRecord(&clients[writerId], ports[writerId], &record)
			time.Sleep(time.Second)

			fmt.Println(record)
		}
	}

	return nil
}

// Определяет на какой Reader (какой порт) отправлять record
func getReaderIDforRecord(r *Record, numWriters uint32) uint32 {
	id := r.ID
	idBytes := make([]byte, 4)

	binary.LittleEndian.PutUint32(idBytes, id)

	crcH := crc32.ChecksumIEEE(idBytes)
	return crcH % numWriters
}

// Посылает запакованный в json record на http сервер с портом port
func postRecord(client *http.Client, port string, record *Record) {
	url := writerUrl + ":" + port

	jsonData, _ := json.Marshal(*record)
	request, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))

	request.Header.Set("Content-Type", "application/json;")
	response, err := client.Do(request)
	if err != nil {
		return
	}

	fmt.Println("response Status:", response.Status)

	defer response.Body.Close()
}

func main() {
	portsPtr := flag.String("w", "6001", "a string")
	filenamePtr := flag.String("f", "amazon-meta.txt.gz", "a string")

	flag.Parse()

	ports := strings.Split(*portsPtr, ",")
	filename := *filenamePtr

	clients := make([]http.Client, len(ports))
	for i := range clients {
		clients[i] = http.Client{}
	}

	fileParser(filename, clients, ports)
}
