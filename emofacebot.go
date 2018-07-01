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

		fileID := getMaxFileID(update.Message.Photo)
		photoURL, err := bot.GetFileDirectURL(fileID)

		if err != nil {
			log.Panic(err)
			return
		}

		faces, err := emo.GetEmotions(photoURL)
		if err != nil {
			// send error
			sendMessage(update.Message.Chat.ID, update.Message.MessageID, err.Error())
			return
		}

		if len(faces) == 0 {
			log.Println("No faces")
			return
		}

		// send emotions
		text := getFacesAsString(faces)

		sendMessage(update.Message.Chat.ID, update.Message.MessageID, text)

		log.Println("Message sent")
	}
}

// send to telegram
func sendMessage(chatID int64, messageID int, message string) {
	msg := tgbotapi.NewMessage(chatID, message)
	msg.ReplyToMessageID = messageID
	bot.Send(msg)
}

// get best with quality image
func getMaxFileID(photos *[]tgbotapi.PhotoSize) string {
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
func getFacesAsString(faces []emotions.Face) string {
	var buffer bytes.Buffer

	isNeedNumeration := len(faces) > 1

	for i, face := range faces {
		if isNeedNumeration {
			buffer.WriteString(fmt.Sprintf("\n#%d:\n", i+1))
		}
		buffer.WriteString(toString(&face.FaceAttributes.Emotion))
	}

	return buffer.String()
}

func toString(emotion *emotions.Emotion) string {
	var buffer bytes.Buffer

	if emotion.Anger > 0 {
		buffer.WriteString(fmt.Sprintf("Anger: %.0f%%\n", emotion.Anger))
	}
	if emotion.Contempt > 0 {
		buffer.WriteString(fmt.Sprintf("Contempt: %.0f%%\n", emotion.Contempt))
	}
	if emotion.Disgust > 0 {
		buffer.WriteString(fmt.Sprintf("Disgust: %.0f%%\n", emotion.Disgust))
	}
	if emotion.Fear > 0 {
		buffer.WriteString(fmt.Sprintf("Fear: %.0f%%\n", emotion.Fear))
	}
	if emotion.Happiness > 0 {
		buffer.WriteString(fmt.Sprintf("Happiness: %.0f%%\n", emotion.Happiness))
	}
	if emotion.Neutral > 0 {
		buffer.WriteString(fmt.Sprintf("Neutral: %.0f%%\n", emotion.Neutral))
	}
	if emotion.Sadness > 0 {
		buffer.WriteString(fmt.Sprintf("Sadness: %.0f%%\n", emotion.Sadness))
	}
	if emotion.Surprise > 0 {
		buffer.WriteString(fmt.Sprintf("Surprise: %.0f%%\n", emotion.Surprise))
	}

	return buffer.String()
}
