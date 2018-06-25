package emofacebot

import (
	"bytes"
	"fmt"
	"log"

	"github.com/aafanasev/ms-emotions-go"
	"gopkg.in/telegram-bot-api.v4"
)

var bot *tgbotapi.BotAPI
var emo *emotions.Emo

func Init(microsoftToken, telegramToken, baseUrl string, debug bool) error {
	var err error

	// create EmoAPI client
	emo = emotions.NewClient(microsoftToken)

	// create a bot
	bot, err = tgbotapi.NewBotAPI(telegramToken)
	if err != nil {
		return err
	}

	log.Printf("[%s] connected", bot.Self.UserName)

	// debug mode
	bot.Debug = debug

	// set webhook
	webhookUrl := baseUrl + bot.Token
	_, err = bot.SetWebhook(tgbotapi.NewWebhook(webhookUrl))
	if err != nil {
		return err
	}

	log.Printf("[%s] set webhook %s", bot.Self.UserName, webhookUrl)

	return err
}

func Handle(update *tgbotapi.Update) {
	if update.Message.Photo != nil {
		log.Println("Photo received")

		fileID := GetMaxFileID(update.Message.Photo)
		photoURL, err := bot.GetFileDirectURL(fileID)

		if err != nil {
			log.Panic(err)
			return
		}

		faces, err := emo.GetEmotions(photoURL)
		if err != nil {
			// send error
			SendMessage(update.Message.Chat.ID, update.Message.MessageID, err.Error())
			return
		}

		if len(faces) == 0 {
			log.Println("No faces")
			return
		}

		// send emotions
		SendMessage(update.Message.Chat.ID, update.Message.MessageID, GetFacesAsString(faces))

		log.Println("Message sent")
	}
}

// send to telegram
func SendMessage(chatID int64, messageID int, message string) {
	msg := tgbotapi.NewMessage(chatID, message)
	msg.ReplyToMessageID = messageID
	bot.Send(msg)
}

// get best with quality image
func GetMaxFileID(photos *[]tgbotapi.PhotoSize) string {
	firstPhoto := (*photos)[0]
	result := firstPhoto.FileID
	width := firstPhoto.Width

	for _, photo := range *photos {
		if width < photo.Width && photo.Width <= 4096 {
			result = photo.FileID
			width = photo.Width
		}
	}

	return result
}

// merge all scores to one message
func GetFacesAsString(faces []emotions.Face) string {
	var buffer bytes.Buffer

	isNeedNumeration := len(faces) > 1

	for i, face := range faces {
		if isNeedNumeration {
			buffer.WriteString(fmt.Sprintf("\n#%d:\n", i+1))
		}
		buffer.WriteString(face.Scores.String())
	}

	return buffer.String()
}
