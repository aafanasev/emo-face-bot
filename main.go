package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Syfaro/telegram-bot-api"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(telegramToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 10

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		log.Println("got message")

		if update.Message.Photo != nil {
			fileID := getMaxFileID(update.Message.Photo)
			photoURL, err := bot.GetFileDirectURL(fileID)

			if err != nil {
				log.Panic(err)
			} else {
				jsonReq := []byte("{\"url\":\"" + photoURL + "\"}")

				req, _ := http.NewRequest("POST", "https://api.projectoxford.ai/emotion/v1.0/recognize", bytes.NewBuffer(jsonReq))
				req.Header.Add("Ocp-Apim-Subscription-Key", emotionsAPIKey)
				req.Header.Add("Content-Type", "application/json")

				client := &http.Client{}
				resp, err := client.Do(req)
				defer resp.Body.Close()

				if err != nil {
					log.Panic(err)
				}

				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					log.Panic(err)
				} else {
					var faces []Face
					err := json.Unmarshal(body, &faces)
					if err != nil {
						log.Panic(err)
					} else {
						for _, face := range faces {
							msg := tgbotapi.NewMessage(update.Message.Chat.ID, face.String())
							msg.ReplyToMessageID = update.Message.MessageID
							bot.Send(msg)
						}
					}
				}
			}
		}
	}
}

func getMaxFileID(photos []tgbotapi.PhotoSize) string {
	result := photos[0].FileID
	width := photos[0].Width

	for _, photo := range photos {
		if photo.Width > width {
			result = photo.FileID
			width = photo.Width
		}
	}

	return result
}
