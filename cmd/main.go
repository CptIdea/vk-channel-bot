package main

import (
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
	"log"
	"vk-channel-bot/bot"
	"vk-channel-bot/cfg"

	"github.com/SevereCloud/vksdk/v2/api"
)

func main() {
	cfg, err := cfg.ParseEnv()
	if err != nil {
		log.Fatal(err)
	}

	vk := api.NewVK(cfg.Token)

	group, err := vk.GroupsGetByID(nil)
	if err != nil {
		log.Fatal(err)
	}

	// Initializing Long Poll
	lp, err := longpoll.NewLongPoll(vk, group[0].ID)
	if err != nil {
		log.Fatal(err)
	}

	ctrl := bot.NewController(vk, lp, log.Default(), cfg.Admins)

	log.Fatalln(ctrl.Start())
}
