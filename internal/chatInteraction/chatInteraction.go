package chatinteraction

import (
	"basic_wallet/internal/db"
	wallet "basic_wallet/internal/manipulations"
	"basic_wallet/internal/state"
	"database/sql"
	"fmt"
	"log"
	"regexp"

	"github.com/ethereum/go-ethereum/ethclient"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func ChatInteractionX(update tgbotapi.Update, dbConn *sql.DB, client *ethclient.Client, bot *tgbotapi.BotAPI, StateOfUser map[string]state.StateOfAddr) (tgbotapi.MessageConfig, error) {
	var msg tgbotapi.MessageConfig
	var err error

	if update.Message.Text == "/start" {
		msg, err = start(update, dbConn, client, StateOfUser)

		if err != nil {
			return tgbotapi.MessageConfig{}, err
		}
	}

	if update.Message.Text == "Показать баланс" {
		msg, err = showBalance(update, dbConn, client)

		if err != nil {
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Не могу забрать балансик :(")
		}
	}

	if update.Message.Text == "Перевести другому человеку эфир" {

		if StateOfUser[update.Message.Chat.UserName] != state.READY_NOT {
			StateOfUser[update.Message.Chat.UserName] = state.READY_NOT
		}

		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Введите ваш адрес:")
	}

	if StateOfUser[update.Message.Chat.UserName] == state.READY_NOT {
		user, err := db.SelectSingleUser(update.Message.Chat.UserName, dbConn)
		addr := update.Message.Text

		if err != nil {
			return tgbotapi.MessageConfig{}, err
		}

		if checkIfValidAddress(addr) {
			wallet.TransferEthereum(client, user.Eth_address, addr, user.Private_key)
		}
	}

	msg.ReplyMarkup = numericKeyboard

	return msg, nil
}

func start(update tgbotapi.Update, dbConn *sql.DB, client *ethclient.Client, StateOfUser map[string]state.StateOfAddr) (tgbotapi.MessageConfig, error) {
	user, err := db.SelectSingleUser(update.Message.Chat.UserName, dbConn)

	if (err != nil && user == db.User{}) {
		log.Println(err)
		addr, privateKey, err := wallet.WalletCreation()

		if err != nil {
			log.Println(err)
			StateOfUser[addr] = state.NOT_NOT
			return tgbotapi.NewMessage(update.Message.Chat.ID, "Мы не смогли создать вам кошелечек! Попробуйте снова отправить /start"), nil
		} else {
			textMsg := fmt.Sprintf("Мы создали вам кошелечек! Его адрес на тестнете: %s", addr)

			db.InsertIntoTable(update.Message.Chat.ID, update.Message.Chat.UserName, addr, privateKey, dbConn)

			return tgbotapi.NewMessage(update.Message.Chat.ID, textMsg), nil
		}

	} else {
		return tgbotapi.NewMessage(update.Message.Chat.ID, "У вас уже есть кошелечек, оберегайте его!"), nil
	}
}

func takeEtherFromOrig(update tgbotapi.Update, dbConn *sql.DB, bot *tgbotapi.BotAPI, addr string) (tgbotapi.MessageConfig, error) {
	_, err := db.SelectSingleUser(update.Message.Chat.UserName, dbConn)

	if err != nil {
		return tgbotapi.MessageConfig{}, nil
	}

	return tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text), nil
}

func showBalance(update tgbotapi.Update, dbConn *sql.DB, client *ethclient.Client) (tgbotapi.MessageConfig, error) {
	usr, _ := db.SelectSingleUser(update.Message.Chat.UserName, dbConn)

	balance, err := wallet.ShowBalance(client, usr.Eth_address)

	if err != nil {
		return tgbotapi.MessageConfig{}, err
	}

	msgTxt := fmt.Sprintf("Ваш балансик составляет: %s ether", balance)
	return tgbotapi.NewMessage(update.Message.Chat.ID, msgTxt), nil
}

var numericKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Перевести другому человеку эфир"),
		tgbotapi.NewKeyboardButton("Показать баланс"),
		tgbotapi.NewKeyboardButton("Взять эфир из основного кошелька"),
	),
)

func checkIfValidAddress(addr string) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")

	return re.MatchString(addr)
}
