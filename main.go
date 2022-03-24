package main

import (
	"log"
	"os"
	"regexp"
	"strings"

	dw "github.com/alessandrovaprio/tele-bot-dwl-yt/downloadHandler"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	ytregex = regexp.MustCompile(`(http:|https:)?\/\/(www\.)?(youtube.com|youtu.be)\/(watch)?(\?v=)?(\S+)?`)
	Version = "0.0.9"
)

func checkMsgIsYoutebVideo(url string) bool {
	return ytregex.MatchString(url)
}
func checkIfDownloadCommand(url string) bool {
	splitted := strings.Split(url, "|")
	return ytregex.MatchString(url) && len(splitted) > 1 && strings.Contains(splitted[0], "download")
}

func doDownloadAndSend(bot *tgbotapi.BotAPI, update tgbotapi.Update, options []string) {
	var chatId int64 = 0

	if update.Message != nil {
		chatId = update.Message.Chat.ID
	}
	if update.CallbackQuery != nil {
		chatId = update.CallbackQuery.Message.Chat.ID
	}

	msgId := 0
	if update.Message != nil {
		msgId = update.Message.MessageID
	}
	if update.CallbackQuery != nil {
		msgId = update.CallbackQuery.Message.MessageID
	}

	msg := tgbotapi.NewMessage(chatId, " Download In progress, may take some minutes ⌛️")
	msg.ReplyToMessageID = msgId
	bot.Send(msg)
	var file tgbotapi.FileBytes
	var err error
	if options[0] == "mp3" {
		file, err = dw.DownloadMp3(options[1])
	} else {
		file, err = dw.DownloadAndConvert(options[1])
	}

	if err != nil {
		sendError(bot, chatId, err)
	} else {
		msg = tgbotapi.NewMessage(chatId, " Download Completed 👍 ✅ 🥳")
		bot.Send(msg)
		if options[0] == "mp3" {
			msgMV := tgbotapi.NewAudio(chatId, file)
			msgMV.ReplyToMessageID = update.CallbackQuery.Message.MessageID
			bot.Send(msgMV)
		} else {
			msgMV := tgbotapi.NewVideo(chatId, file)
			msgMV.ReplyToMessageID = update.CallbackQuery.Message.MessageID
			bot.Send(msgMV)
		}
	}
}

func sendError(bot *tgbotapi.BotAPI, chatId int64, err error) {
	msg := tgbotapi.NewMessage(chatId, " 👎 Error "+err.Error())
	msg.ReplyToMessageID = int(chatId)
	bot.Send(msg)
}
func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_API_KEY"))
	if err != nil {
		log.Panic(err)
	}

	// bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 10

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message

			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			if checkMsgIsYoutebVideo(update.Message.Text) {
				log.Printf("%s", "choose a format")
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "choose a format")
				msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("📽 Video (.mpeg)", "mpeg|"+update.Message.Text),
						tgbotapi.NewInlineKeyboardButtonData("🎵 Audio (.mp3)", "mp3|"+update.Message.Text),
					),
				)
				bot.Send(msg)

			} else {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, " No youtube url 👎")
				bot.Send(msg)
			}
		}
		// check if it's a response from button
		if update.CallbackQuery != nil {
			splitted := strings.Split(update.CallbackQuery.Data, "|")
			if len(splitted) > 1 {
				go doDownloadAndSend(bot, update, splitted)
			}
			log.Printf(update.CallbackQuery.Data)
		}
	}
}
