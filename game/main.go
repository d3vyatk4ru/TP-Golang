package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	badCommand         = "неизвестная команда"
	badLoot            = "нет такого"
	badBag             = "некуда класть"
	tmplLoot           = "предмет добавлен в инвентарь: "
	tmplApply          = "вы надели: "
	badLootInInventory = "нет предмета в инвентаре - "
	badToApply         = "не к чему применить"
	badWay             = "нет пути в "
	close              = " закрыта"
	open               = " открыта"
)

var player Player
var room = make(map[string]Room, len(getData("rooms.txt")))
var loot = make([]string, len(getData("rooms.txt")))

type Item struct {
	KeyTo string
	flag  bool
}

type Room struct {
	Name        string
	Description map[string][]string
	exit        map[string]bool
	Loot        map[string]bool
	Need        map[string]bool
	Interaction map[string]bool
	State       int
}

type Player struct {
	CurrentRoom Room
	Dressed     map[string]Item
	Loot        map[string]Item
}

func getData(fileName string) []string {

	f, err := os.Open(fileName)

	if err != nil {
		fmt.Println("ooops!!!!")
	}

	scanner := bufio.NewScanner(f)

	var data []string

	for scanner.Scan() {

		data = append(data, scanner.Text())
	}

	return data
}

func getItem(playerItem map[string]Item, command []string) string {

	for key, val := range playerItem {

		_, found := room[player.CurrentRoom.Name].Loot[key]

		if len(command) == 2 && key == command[1] && found {

			if command[1] != "рюкзак" && !player.Dressed["рюкзак"].flag {
				return badBag
			}

			if val.flag {
				return badLoot
			}

			if command[0] == "надеть" {
				player.Dressed[key] = Item{player.Loot[key].KeyTo, true}
			} else {
				player.Loot[key] = Item{player.Loot[key].KeyTo, true}
			}

			_, found = room[player.CurrentRoom.Name].Loot[command[1]]

			if found {
				room[player.CurrentRoom.Name].Loot[command[1]] = false
			}

			if command[0] == "надеть" {
				player.CurrentRoom.State++
				return tmplApply + key
			}

			player.CurrentRoom.State++
			return tmplLoot + key
		}
	}

	return badLoot
}

func applyItem(command []string) string {

	_, foundLoot := player.Loot[command[1]]

	// есть ли вещь в интвенторе
	if !player.Loot[command[1]].flag {
		return badLootInInventory + command[1]
	}

	// есть ли вещь в комнате
	if !foundLoot {
		return badLoot
	}

	_, foundToApplying := player.CurrentRoom.Interaction[command[2]]

	// есть ли вещь в комнате
	if !foundToApplying {
		return badToApply
	}

	// если уже взяли вещь
	for !player.CurrentRoom.Interaction[command[2]] {
		return badToApply
	}

	player.CurrentRoom.Interaction[command[2]] = false
	player.CurrentRoom.exit[player.Loot[command[1]].KeyTo] = true

	return command[2] + open
}

func checkKitchen(player Player) bool {

	if player.CurrentRoom.Name == "кухня" {
		for _, val := range player.Loot {
			if !val.flag {
				return false
			}
		}

		return true
	}

	return false
}

func (player *Player) do(command []string) string {

	switch command[0] {

	case "осмотреться":

		if checkKitchen(*player) {
			player.CurrentRoom.State++
			return player.CurrentRoom.Description[command[0]][player.CurrentRoom.State]
		}

		return player.CurrentRoom.Description[command[0]][player.CurrentRoom.State]

	case "идти":
		for key1, val1 := range player.CurrentRoom.exit {

			if len(command) == 2 && key1 == command[1] {

				if val1 {

					player.CurrentRoom = room[key1]
					return player.CurrentRoom.Description[command[0]][player.CurrentRoom.State]

				} else {

					for key2, val2 := range player.CurrentRoom.Need {
						if !val2 {
							return key2 + close
						}
					}

					for key2, val2 := range player.CurrentRoom.Interaction {
						if val2 {
							return key2 + close
						}
					}
				}
			}
		}

		return badWay + command[1]

	case "надеть":
		return getItem(player.Dressed, command)

	case "взять":
		return getItem(player.Loot, command)

	case "применить":
		return applyItem(command)
	}

	return badCommand
}

func main() {

	initGame()

	fmt.Println(handleCommand("осмотреться"))
	fmt.Println(handleCommand("завтракать"))
	fmt.Println(handleCommand("идти комната"))
	fmt.Println(handleCommand("идти коридор"))
	fmt.Println(handleCommand("применить ключи дверь"))
	fmt.Println(handleCommand("идти комната"))
	fmt.Println(handleCommand("осмотреться"))
	fmt.Println(handleCommand("надеть рюкзак"))
	fmt.Println(handleCommand("осмотреться"))
	fmt.Println(handleCommand("взять ключи"))
	fmt.Println(handleCommand("взять телефон"))
	fmt.Println(handleCommand("взять ключи"))
	fmt.Println(handleCommand("осмотреться"))
	fmt.Println(handleCommand("взять конспекты"))
	fmt.Println(handleCommand("осмотреться"))
	fmt.Println(handleCommand("идти коридор"))
	fmt.Println(handleCommand("идти кухня"))
	fmt.Println(handleCommand("осмотреться"))
	fmt.Println(handleCommand("идти коридор"))
	fmt.Println(handleCommand("идти улица"))
	fmt.Println(handleCommand("применить ключи дверь"))
	fmt.Println(handleCommand("применить телефон шкаф"))
	fmt.Println(handleCommand("применить ключи шкаф"))
	fmt.Println(handleCommand("идти улица"))
}

func initGame() {

	room["кухня"] = Room{
		Name: "кухня",
		Description: map[string][]string{"осмотреться": []string{"ты находишься на кухне, на столе: чай, надо собрать рюкзак и идти в универ. можно пройти - коридор",
			"ты находишься на кухне, на столе: чай, надо идти в универ. можно пройти - коридор"},
			"идти": []string{"кухня, ничего интересного. можно пройти - коридор"}},
		exit: map[string]bool{"коридор": true}}

	room["коридор"] = Room{
		Name: "коридор",
		Description: map[string][]string{"осмотреться": []string{""},
			"идти": []string{"ничего интересного. можно пройти - кухня, комната, улица"}},
		exit:        map[string]bool{"кухня": true, "комната": true, "улица": false},
		Interaction: map[string]bool{"дверь": true}}

	room["комната"] = Room{
		Name: "комната",
		Description: map[string][]string{"осмотреться": []string{"на столе: ключи, конспекты, на стуле: рюкзак. можно пройти - коридор",
			"на столе: ключи, конспекты. можно пройти - коридор",
			"на столе: конспекты. можно пройти - коридор",
			"пустая комната. можно пройти - коридор"},
			"идти": []string{"ты в своей комнате. можно пройти - коридор"}},
		exit: map[string]bool{"коридор": true},
		Loot: map[string]bool{"ключи": true, "конспекты": true, "рюкзак": true}}

	room["улица"] = Room{
		Name: "улица",
		Description: map[string][]string{"осмотреться": []string{""},
			"идти": []string{"на улице весна. можно пройти - домой"}},
		Need: map[string]bool{"ключи": false},
	}

	player.CurrentRoom = room["кухня"]

	player.Loot = make(map[string]Item)

	player.Loot["ключи"] = Item{"улица", false}
	player.Loot["конспекты"] = Item{"", false}

	player.Dressed = make(map[string]Item)
	player.Dressed["рюкзак"] = Item{"", false}

	for i, val := range getData("loot.txt") {
		loot[i] = val
	}

}

func handleCommand(command string) string {

	splitCommand := strings.Split(command, " ")

	return player.do(splitCommand)
}
