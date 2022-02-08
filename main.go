package main

import (
	"context"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/events"
)

func inArray(id int, array []int) bool {
	for _, i := range array {
		if i == id {
			return true
		}
	}
	return false
}

func main() {
	var admins []int
	adminsString := strings.Split(os.Getenv("ADMINS"), ",")
	for _, s := range adminsString {
		id, err := strconv.Atoi(s)
		if err != nil {
			log.Fatalf("ошибка парсинга id админа(%s): %s ", s, err)
		}
		admins = append(
			admins,
			id,
		)
	}

	token := os.Getenv("TOKEN")
	if token == "" {
		log.Fatalln("TOKEN не должен быть пустым")
	}
	vk := api.NewVK(token)

	// get information about the group
	group, err := vk.GroupsGetByID(nil)
	if err != nil {
		log.Fatal(err)
	}

	// Initializing Long Poll
	lp, err := longpoll.NewLongPoll(vk, group[0].ID)
	if err != nil {
		log.Fatal(err)
	}

	// New message event
	lp.MessageNew(func(_ context.Context, obj events.MessageNewObject) {
		log.Printf("%d: %s (%d:%d)", obj.Message.FromID, obj.Message.Text, obj.Message.PeerID, obj.Message.ConversationMessageID)

		if !inArray(obj.Message.FromID, admins) {
			_, err := vk.MessagesDelete(params.NewMessagesDeleteBuilder().DeleteForAll(true).ConversationMessageIDs([]int{obj.Message.ConversationMessageID}).PeerID(obj.Message.PeerID).Params)
			if err != nil {
				if !strings.Contains(err.Error(), "Access denied: message can not be deleted (admin message)") {
					log.Println(err)
				}
				return
			}
		}
	})

	// Run Bots Long Poll
	log.Println("Start Long Poll")
	log.Println("Бот запущен со списком админов:", admins)
	if err := lp.Run(); err != nil {
		log.Fatal(err)
	}
}
