package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

var player Player
var room []Room

func getData(fileName string) []string {

	f, err := os.Open(fileName)

	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(f)

	var data []string

	for scanner.Scan() {

		data = append(data, scanner.Text())
	}

	return data
}

func getString(fileName string, num_str int) string {

	f, err := os.Open(fileName)

	if err != nil {
		fmt.Println("ooops!!!!")
	}

	scanner := bufio.NewScanner(f)

	var data string

	for i := 0; i <= num_str; i++ {

		scanner.Scan()

		if i == num_str {
			data = scanner.Text()
		}
	}

	return data
}

type Room struct {
	Title       string
	Id          int
	Description string
	exit        []string
	Invetory    map[string]bool
}

type Player struct {
	CurrentRoom string
	Invetory    map[string]bool
}

func (player *Player) do(userAction string) {

}

func main() {
	/*
		в этой функции можно ничего не писать
		но тогда у вас не будет работать через go run main.go
		очень круто будет сделать построчный ввод команд тут, хотя это и не требуется по заданию
	*/

	// initGame()

	handleCommand("применить ключи шкаф")

}

/*
	эта функция инициализирует игровой мир - все команты
	если что-то было - оно корректно перезатирается
*/

func initGame() {

	rooms := getData("rooms.txt")
	player.CurrentRoom = rooms[0]

	room = make([]Room, 2*len(rooms), 2*len(rooms))

	for i, val := range rooms {
		room[i].Title = val
		room[i].Id = i
		room[i].Description = getData("initDescription.txt")[i]
		room[i].exit = strings.Split(getString("exit.txt", i), " ")
		room[i].Invetory = make(map[string]bool)

		if str := getString("inventory.txt", i); str != "" {
			for _, val := range strings.Split(str, " ") {
				room[i].Invetory[val] = true
			}
		}
	}

	player.Invetory = make(map[string]bool)
	for _, val := range getData("things.txt") {
		player.Invetory[val] = false
	}
}

func handleCommand(command string) string {
	/*
		данная функция принимает команду от "пользователя"
		и наверняка вызывает какой-то другой метод или функцию у "мира" - списка комнат
	*/

	splt := strings.Split(command, " ")

	return "not implemented"
}
