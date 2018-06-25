package emofacebot

import (
	"bytes"
	"fmt"
	"log"
	"net/http"

	"github.com/jokuskay/ms-emotions-go"
	"gopkg.in/telegram-bot-api.v4"
)

func main() {
	// create EmoAPI client
	emo := emotions.NewClient("")

	// create Telegram bot client
	bot, err := tgbotapi.NewBotAPI("")
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	_, err = bot.SetWebhook(tgbotapi.NewWebhookWithCert("[SITE URL]"+bot.Token, "[CERT.PEM]"))
	if err != nil {
		log.Fatal(err)
	}

	updates, _ := bot.ListenForWebhook("/" + bot.Token)
	go http.ListenAndServeTLS(":[PORT]", "[CERT.PEM]", "[KEY.PEM]", nil)

	for update := range updates {
		if update.Message.Photo != nil {
			log.Println("Photo received")

			fileID := getMaxFileID(update.Message.Photo)
			photoURL, err := bot.GetFileDirectURL(fileID)

			if err != nil {
				log.Panic(err)
				continue
			}

			faces, err := emo.GetEmotions(photoURL)
			if err != nil {
				// send error
				sendMessage(bot, update.Message.Chat.ID, update.Message.MessageID, err.Error())
				continue
			}

			if len(faces) == 0 {
				log.Println("No faces")
				continue
			}

			// send emotions
			sendMessage(bot, update.Message.Chat.ID, update.Message.MessageID, getFacesAsString(faces))

			log.Println("Message sent")
		}
	}
}

// send to telegram
func sendMessage(bot *tgbotapi.BotAPI, chatID int, messageID int, message string) {
	msg := tgbotapi.NewMessage(chatID, message)
	msg.ReplyToMessageID = messageID
	bot.Send(msg)
}

// get best with quality image
func getMaxFileID(photos []tgbotapi.PhotoSize) string {
	result := photos[0].FileID
	width := photos[0].Width

	for _, photo := range photos {
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
		buffer.WriteString(face.Scores.String())
	}

	return buffer.String()
}
