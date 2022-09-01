package main

import (
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	tgbot, err := NewTgBot()
	if err != nil {
		panic(err)
	}

	for update := range tgbot.GetUpdates() {
		if update.Message.Document != nil && update.Message.Document.MimeType == "image/heic" {
			file, err := tgbot.GetDocument(update.Message.Document.FileID)
			if err != nil {
				log.Printf("GetDocument error: %v\n", err)
				continue
			}

			jpgFile, err := HeicToJpg(file)
			if err != nil {
				log.Printf("HeicToJpg error: %v\n", err)
				continue
			}
			jpgName := strings.TrimRight(update.Message.Document.FileName, "heic")
			jpgName = strings.TrimRight(jpgName, "HEIC")
			jpgName = fmt.Sprintf("%s.jpg", jpgName)

			imageToSend := tgbotapi.NewDocument(update.Message.Chat.ID, tgbotapi.FileBytes{
				Name:  jpgName,
				Bytes: jpgFile,
			})

			_, err = tgbot.BotApi.Send(imageToSend)
			if err != nil {
				log.Printf("Send error: %v\n", err)
				continue
			}
		}
	}

}
