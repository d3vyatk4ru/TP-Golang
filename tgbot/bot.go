package main

// сюда писать код

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	tgbotapi "github.com/skinass/telegram-bot-api/v5"
)

type UserInfo struct {
	text     string
	nickname string
	taskID   []int
}

type TaskInfo struct {
	// задача
	text string
	// автор задачи
	authorID int64
	// ID владельца
	ownerID int64
	// никнейм владельца
	authorNickName string
	// никнейм владельца
	ownerNickName string
	// назначена на коко-либо
	assign bool
}

var (
	// @BotFather в телеграме даст вам это
	BotToken = "5693325206:AAFnq-ayD4rq8nz9959-284QtBo2L2cFtkc"

	// урл выдаст вам игрок или хероку
	WebhookURL = "https://3b7f-95-165-106-200.eu.ngrok.io"

	// порт
	port = "8081"

	// храним задачи пользвателей по их ID
	tasksByUserID = make(map[int64]UserInfo)

	// задачи
	nTask = make([]int, 0, 1)

	// задача по её ID
	tasksInfoByTaskID = make(map[int64]TaskInfo)

	// текущая задача
	taskID = 0

	// regex for assign_*
	assign = regexp.MustCompile(`assign_\d+`)

	// regex for unassign_*
	unassign = regexp.MustCompile(`unassign_\d+`)

	// regex for resolve_*
	resolve = regexp.MustCompile(`resolve_\d+`)
)

// удаление жлемента по значению
func remove(slice []int, item int) []int {
	for i, other := range slice {
		if other == item {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

// проверка наличия элемента
func contains(slice []int, target int) bool {
	for _, element := range slice {
		if element == target {
			return true
		}
	}
	return false
}

func taskHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) error {

	if len(tasksInfoByTaskID) == 0 {
		_, err := bot.Send(tgbotapi.NewMessage(
			update.Message.From.ID,
			"Нет задач",
		))

		if err != nil {
			fmt.Println("Can't send message to user, ", err.Error())
			return err
		}

	} else {

		var answerString string
		var idx int
		var resp string

		for _, task := range tasksInfoByTaskID {

			if !tasksInfoByTaskID[int64(nTask[idx])].assign {

				if answerString != "" {
					answerString += "\n\n"
				}

				answerString += fmt.Sprintf(
					"%d. %s by @%s\n/assign_%d", nTask[idx],
					tasksInfoByTaskID[int64(nTask[idx])].text,
					tasksInfoByTaskID[int64(nTask[idx])].authorNickName, nTask[idx],
				)
			} else {

				if task.ownerID == update.Message.From.ID {
					resp = "я"
				} else {
					resp = fmt.Sprintf("@%s",
						tasksInfoByTaskID[int64(nTask[idx])].ownerNickName)
				}

				answerString += fmt.Sprintf(
					"%d. %s by @%s\nassignee: %s",
					nTask[idx],
					task.text,
					tasksInfoByTaskID[int64(nTask[idx])].authorNickName,
					resp,
				)

				if update.Message.From.ID == tasksInfoByTaskID[int64(nTask[idx])].ownerID {
					answerString += fmt.Sprintf("\n/unassign_%d /resolve_%d",
						nTask[idx], nTask[idx])
				}
			}

			idx++
		}

		_, err := bot.Send(tgbotapi.NewMessage(
			update.Message.From.ID,
			answerString,
		))

		if err != nil {
			fmt.Println("Error")
			return err
		}
	}

	return nil
}

func newHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) error {

	if update.Message.CommandArguments() == "" {
		_, err := bot.Send(tgbotapi.NewMessage(
			update.Message.From.ID,
			"Нет задачи!",
		))

		if err != nil {
			fmt.Println("Error")
		}

		return nil
	}

	args := update.Message.CommandArguments()

	taskID++
	nTask = append(nTask, taskID)
	tasksInfoByTaskID[int64(taskID)] = TaskInfo{
		text:           args,
		ownerID:        -1,
		ownerNickName:  "",
		authorID:       update.Message.From.ID,
		authorNickName: update.Message.From.UserName,
	}

	tasksByUserID[update.Message.From.ID] = UserInfo{
		args, update.Message.From.UserName, []int{taskID},
	}

	_, err := bot.Send(tgbotapi.NewMessage(
		update.Message.From.ID,
		fmt.Sprintf(
			"Задача \"%s\" создана, id=%d", args, taskID,
		),
	))

	if err != nil {
		fmt.Println("Can't send message to user, ", err.Error())
		return err
	}

	return nil
}

func unassignHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update, cmd string) error {

	ID, err := strconv.Atoi(strings.Split(cmd, "_")[1])

	if err != nil {
		fmt.Println("Can't cast ID to int, ", err.Error())
		return err
	}

	if !contains(nTask, ID) {
		_, err2 := bot.Send(tgbotapi.NewMessage(
			update.Message.From.ID,
			"Задач с таким номером нет!",
		))

		if err2 != nil {
			fmt.Println("Error")
			return err2
		}

		return nil
	}

	if update.Message.From.ID != tasksInfoByTaskID[int64(ID)].ownerID {
		_, err2 := bot.Send(tgbotapi.NewMessage(
			update.Message.From.ID,
			"Задача не на вас",
		))

		if err2 != nil {
			fmt.Println("Error")
			return err
		}

	} else {
		tasksInfoByTaskID[int64(ID)] = TaskInfo{
			text:           tasksInfoByTaskID[int64(ID)].text,
			ownerID:        -1,
			assign:         false,
			ownerNickName:  "",
			authorID:       tasksInfoByTaskID[int64(ID)].authorID,
			authorNickName: tasksInfoByTaskID[int64(ID)].authorNickName,
		}

		_, err = bot.Send(tgbotapi.NewMessage(
			update.Message.From.ID,
			"Принято",
		))

		if err != nil {
			fmt.Println("Bad sent msg for user", err.Error())
			return err
		}

		_, err := bot.Send(tgbotapi.NewMessage(
			tasksInfoByTaskID[int64(ID)].authorID,
			fmt.Sprintf("Задача \"%s\" осталась без исполнителя", tasksInfoByTaskID[int64(ID)].text),
		))

		if err != nil {
			fmt.Println("Error")
			return err
		}
	}

	return nil
}

func assignHander(bot *tgbotapi.BotAPI, update tgbotapi.Update, cmd string) error {

	var oldOwnerID int64

	ID, err := strconv.Atoi(strings.Split(cmd, "_")[1])

	if err != nil {
		fmt.Println("Can't cast ID to int, ", err.Error())
		return err
	}

	if !contains(nTask, ID) {

		_, err2 := bot.Send(tgbotapi.NewMessage(
			update.Message.From.ID,
			"Задач с таким номером нет!",
		))

		if err2 != nil {
			fmt.Println("Error")
			return err2
		}

		return nil
	}

	if tasksInfoByTaskID[int64(ID)].ownerID == -1 {
		oldOwnerID = tasksInfoByTaskID[int64(ID)].authorID
	} else {
		oldOwnerID = tasksInfoByTaskID[int64(ID)].ownerID
	}

	tasksInfoByTaskID[int64(ID)] = TaskInfo{
		text:           tasksInfoByTaskID[int64(ID)].text,
		ownerID:        update.Message.From.ID,
		assign:         true,
		ownerNickName:  update.Message.From.UserName,
		authorNickName: tasksInfoByTaskID[int64(ID)].authorNickName,
		authorID:       tasksInfoByTaskID[int64(ID)].authorID,
	}

	if oldOwnerID != tasksInfoByTaskID[int64(ID)].ownerID {

		_, err2 := bot.Send(tgbotapi.NewMessage(
			oldOwnerID,
			fmt.Sprintf("Задача \"%s\" назначена на @%s",
				tasksInfoByTaskID[int64(ID)].text,
				update.Message.From.UserName),
		))

		if err2 != nil {
			fmt.Println("Error")
			return nil
		}
	}

	_, err = bot.Send(tgbotapi.NewMessage(
		update.Message.From.ID,
		fmt.Sprintf("Задача \"%s\" назначена на вас",
			tasksInfoByTaskID[int64(ID)].text),
	))

	if err != nil {
		fmt.Println("Error")
		return err
	}

	return nil
}

func resolveHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update, cmd string) error {
	ID, err := strconv.Atoi(strings.Split(cmd, "_")[1])

	if err != nil {
		fmt.Println("Can't cast ID to int, ", err.Error())
		return err
	}

	if !contains(nTask, ID) {
		_, err := bot.Send(tgbotapi.NewMessage(
			update.Message.From.ID,
			fmt.Sprintf("Задачи %d не существует", ID),
		))

		if err != nil {
			fmt.Println("Error")
			return err
		}

		return nil
	}

	if tasksInfoByTaskID[int64(ID)].ownerID == update.Message.From.ID {

		authorID := tasksInfoByTaskID[int64(ID)].authorID
		ownerID := tasksInfoByTaskID[int64(ID)].ownerID

		answerToOwner := fmt.Sprintf("Задача \"%s\" выполнена",
			tasksInfoByTaskID[int64(ID)].text)

		answerToAuthor := fmt.Sprintf("Задача \"%s\" выполнена @%s",
			tasksInfoByTaskID[int64(ID)].text, tasksInfoByTaskID[int64(ID)].ownerNickName)

		_, err := bot.Send(tgbotapi.NewMessage(
			ownerID,
			answerToOwner,
		))

		if err != nil {
			fmt.Println("Error")
			return err
		}

		if update.Message.From.ID != authorID {

			_, err := bot.Send(tgbotapi.NewMessage(
				authorID,
				answerToAuthor,
			))

			if err != nil {
				fmt.Println("Error")
				return err
			}
		}

		nTask = remove(nTask, ID)
		delete(tasksInfoByTaskID, int64(ID))
	}

	return nil
}

func myHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) error {

	var answer string

	for taskID, val := range tasksInfoByTaskID {

		if val.ownerID == update.Message.From.ID {

			if answer != "" {
				answer += "\n\n"
			}

			answer += fmt.Sprintf("%d. %s by @%s\n/unassign_%d /resolve_%d",
				taskID,
				val.text,
				val.ownerNickName,
				taskID,
				taskID)
		}
	}

	if answer == "" {
		answer = "Задач нет!"
	}

	_, err := bot.Send(tgbotapi.NewMessage(
		update.Message.From.ID,
		answer,
	))

	if err != nil {
		fmt.Println("Error")
		return err
	}

	return nil
}

func ownerHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) error {

	var answer string

	for taskID, val := range tasksInfoByTaskID {

		if val.authorID == update.Message.From.ID {

			if answer != "" {
				answer += "\n\n"
			}

			answer += fmt.Sprintf("%d. %s by @%s\n/assign_%d",
				taskID,
				val.text,
				val.authorNickName,
				taskID,
			)
		}
	}

	if answer == "" {
		answer = "Задач нет!"
	}

	_, err := bot.Send(tgbotapi.NewMessage(
		update.Message.From.ID,
		answer,
	))

	if err != nil {
		fmt.Println("Error")
		return err
	}

	return nil
}

func defaultHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
	_, err := bot.Send(tgbotapi.NewMessage(
		update.Message.From.ID,
		fmt.Sprintf(
			"Набор команд:\n1. /tasks\n2. /new <task_name>\n"+
				"3. /assign_*\n4. /unassign_*\n5./resolve_*\n6. /my\n7. /owner",
		),
	))

	if err != nil {
		fmt.Println("Error")
		return err
	}

	return nil
}

func startTaskBot(ctx context.Context) error {

	bot, err := tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		log.Fatalf("NewBotAPI failed: %s", err)
	}

	bot.Debug = true

	wh, err := tgbotapi.NewWebhook(WebhookURL)
	if err != nil {
		log.Fatalf("NewWebhook failed: %s", err)
	}

	_, err = bot.Request(wh)
	if err != nil {
		log.Fatalf("SetWebhook failed: %s", err)
	}

	updates := bot.ListenForWebhook("/")

	go func() {
		log.Fatalln("http err:", http.ListenAndServe(":"+port, nil))
	}()

	fmt.Println("start listen :" + port)

	for update := range updates {

		if update.Message == nil {
			continue
		}

		switch cmd := update.Message.Command(); {

		case cmd == "tasks":
			err = taskHandler(bot, update)
			if err != nil {
				fmt.Println("Erorr", err.Error())
			}

		case cmd == "new":
			err = newHandler(bot, update)
			if err != nil {
				fmt.Println("Erorr", err.Error())
			}

		case unassign.MatchString(cmd):
			err = unassignHandler(bot, update, cmd)
			if err != nil {
				fmt.Println("Erorr", err.Error())
			}

		case assign.MatchString(cmd):
			err = assignHander(bot, update, cmd)
			if err != nil {
				fmt.Println("Erorr", err.Error())
			}

		case resolve.MatchString(cmd):
			err = resolveHandler(bot, update, cmd)
			if err != nil {
				fmt.Println("Erorr", err.Error())
			}

		case cmd == "my":
			err = myHandler(bot, update)
			if err != nil {
				fmt.Println("Erorr", err.Error())
			}

		case cmd == "owner":
			err = ownerHandler(bot, update)
			if err != nil {
				fmt.Println("Erorr", err.Error())
			}

		default:
			err = defaultHandler(bot, update)
			if err != nil {
				fmt.Println("Erorr", err.Error())
			}
		}
	}
	return nil
}

func main() {
	err := startTaskBot(context.Background())
	if err != nil {
		panic(err)
	}
}
