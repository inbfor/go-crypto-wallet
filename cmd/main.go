package main

import (
	"flag"
	"log"

	chatinteraction "basic_wallet/internal/chatInteraction"
	db "basic_wallet/internal/db"
	"basic_wallet/internal/state"

	"github.com/ethereum/go-ethereum/ethclient"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {

	apiKey := flag.String("apiKey", "", "Api Key for bot From BotFather")
	flag.Parse()

	bot, err := tgbotapi.NewBotAPI(*apiKey)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	StateOfUser := map[string]state.StateOfAddr{}
	client, _ := ethclient.Dial("https://rpc.teku-geth-001.srv.holesky.ethpandaops.io")
	dbConn, err := db.Connect("users.db")

	if err != nil {
		log.Println(err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {

		msg, err := chatinteraction.ChatInteractionX(update, dbConn, client, bot, StateOfUser)

		if err != nil {
			log.Println(err)
		}

		bot.Send(msg)

	}
}
