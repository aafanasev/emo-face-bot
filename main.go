package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Syfaro/telegram-bot-api"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(telegramToken)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 10

	updates, _ := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message.Photo != nil {
			log.Println("Got photo")

			fileID := getMaxFileID(update.Message.Photo)
			photoURL, err := bot.GetFileDirectURL(fileID)

			if err != nil {
				log.Panic(err)
				continue
			}

			jsonReq := []byte("{\"url\":\"" + photoURL + "\"}")

			req, _ := http.NewRequest("POST", "https://api.projectoxford.ai/emotion/v1.0/recognize", bytes.NewBuffer(jsonReq))
			req.Header.Add("Ocp-Apim-Subscription-Key", emotionsAPIKey)
			req.Header.Add("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			defer resp.Body.Close()

			if err != nil {
				log.Panic(err)
				continue
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Panic(err)
				continue
			}

			var faces []Face
			err = json.Unmarshal(body, &faces)
			if err != nil {
				log.Panic(err)
				continue
			}

			if len(faces) == 0 {
				log.Println("No faces")
				continue
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, getFacesAsString(faces))
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)

			log.Println("Message sent")
		}
	}
}

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

func getFacesAsString(faces []Face) string {
	var buffer bytes.Buffer

	isNeedNumeration := len(faces) > 1

	for i, face := range faces {
		if isNeedNumeration {
			buffer.WriteString(fmt.Sprintf("\n#%d:\n", i+1))
		}
		buffer.WriteString(face.String())
	}

	return buffer.String()
}

func round(val float32) int {
	if val < 0 {
		return int(val - 0.5)
	}
	return int(val + 0.5)
}
