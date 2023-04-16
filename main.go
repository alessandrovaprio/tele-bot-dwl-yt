package main

import (
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	dw "github.com/alessandrovaprio/tele-bot-dwl-yt/downloadHandler"
	"github.com/alessandrovaprio/tele-bot-dwl-yt/fileHandler"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	ytregex = regexp.MustCompile(`(http:|https:)?\/\/(www\.)?(youtube.com|youtu.be)\/(watch)?(\?v=)?(\S+)?`)
	Version = "0.0.9"
)
var video_urls = make(map[string]string)
var video_index = 0

func checkMsgIsYoutebVideo(url string) bool {
	return ytregex.MatchString(url)
}
func checkIfDownloadCommand(url string) bool {
	splitted := strings.Split(url, "|")
	return ytregex.MatchString(url) && len(splitted) > 1 && strings.Contains(splitted[0], "download")
}

func doDownloadAndSend(bot *tgbotapi.BotAPI, update tgbotapi.Update, option string, video_url_key string) {
	var chatId int64 = 0
	// remove the key-value pair when finish
	defer delete(video_urls, video_url_key)

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
	if option == "mp3" {
		file, err = dw.DownloadMp3(option)
	} else {
		file, err = dw.DownloadAndConvert(option)
	}

	if err != nil {
		sendError(bot, chatId, err)
	} else {
		msg = tgbotapi.NewMessage(chatId, " Download Completed 👍 ✅ 🥳")
		bot.Send(msg)
		if option == "mp3" {
			msgMV := tgbotapi.NewAudio(chatId, file)
			msgMV.ReplyToMessageID = update.CallbackQuery.Message.MessageID
			bot.Send(msgMV)
			// in any case I delete the mp3 file
			defer fileHandler.DeleteFile(file.Name)
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
			video_index = 0
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			if checkMsgIsYoutebVideo(update.Message.Text) {
				for {
					_, found := video_urls[strconv.Itoa(video_index)]
					if found {
						video_index++
					}
					break
				}
				video_urls[strconv.Itoa(video_index)] = update.Message.Text
				log.Printf("%s", "choose a format")
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "choose a format")
				msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("📽 Video (.mpeg)", "mpeg|"+strconv.Itoa(video_index)),
						tgbotapi.NewInlineKeyboardButtonData("🎵 Audio (.mp3)", "mp3|"+strconv.Itoa(video_index)),
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
				value, found := video_urls[splitted[1]]
				if found {
					go doDownloadAndSend(bot, update, value, splitted[1])
				}
			}
			log.Printf(update.CallbackQuery.Data)
		}
	}
}
