package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/caarlos0/env/v6"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TgBot struct {
	ApiKey     string `env:"BOT_API_KEY,required"`
	BotApi     *tgbotapi.BotAPI
	HttpClient *http.Client
}

func NewTgBot() (*TgBot, error) {
	tgBot := &TgBot{}
	err := tgBot.ParseEnvs()

	bot, err := tgbotapi.NewBotAPI(tgBot.ApiKey)
	if err != nil {
		return nil, err
	}

	tgBot.BotApi = bot
	tgBot.HttpClient = &http.Client{
		Timeout: 10 * time.Second,
	}

	return tgBot, nil
}

func (t *TgBot) ParseEnvs() error {
	err := env.Parse(t)
	if err != nil {
		return fmt.Errorf("env parsing error: %v", err)
	}
	return nil
}

func (t *TgBot) GetUpdates() tgbotapi.UpdatesChannel {
	return t.BotApi.GetUpdatesChan(tgbotapi.UpdateConfig{
		Offset:         0,
		Limit:          1,
		Timeout:        1,
		AllowedUpdates: nil,
	})
}

func (t *TgBot) GetDocument(fileId string) ([]byte, error) {
	filePath, err := t.getFilePath(fileId)
	if err != nil {
		return nil, err
	}

	file, err := t.getFile(filePath)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (t *TgBot) getFilePath(fileId string) (string, error) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/getFile?file_id=%s", t.BotApi.Token, fileId)
	resp, err := t.HttpClient.Get(url)
	defer resp.Body.Close()
	if err != nil {
		return "", err
	}
	type getFileResponse struct {
		Ok     bool `json:"ok"`
		Result struct {
			FileId       string `json:"file_id"`
			FileUniqueId string `json:"file_unique_id"`
			FileSize     int    `json:"file_size"`
			FilePath     string `json:"file_path"`
		} `json:"result"`
	}

	response := &getFileResponse{}
	err = json.NewDecoder(resp.Body).Decode(response)
	if err != nil {
		return "", err
	}
	return response.Result.FilePath, nil
}

func (t *TgBot) getFile(filePath string) ([]byte, error) {
	url := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", t.BotApi.Token, filePath)
	resp, err := t.HttpClient.Get(url)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	byteFile, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return byteFile, nil
}

func (t *TgBot) Heic2JpgDoc(message *tgbotapi.Message) {
	file, err := t.GetDocument(message.Document.FileID)
	if err != nil {
		log.Printf("GetDocument error: %v\n", err)
		return
	}

	jpgFile, err := HeicToJpg(file)
	if err != nil {
		log.Printf("HeicToJpg error: %v\n", err)
		return
	}
	jpgName := strings.TrimRight(message.Document.FileName, "heic")
	jpgName = strings.TrimRight(jpgName, "HEIC")
	jpgName = fmt.Sprintf("%s.jpg", jpgName)

	imageToSend := tgbotapi.NewDocument(message.Chat.ID, tgbotapi.FileBytes{
		Name:  jpgName,
		Bytes: jpgFile,
	})

	_, err = t.BotApi.Send(imageToSend)
	if err != nil {
		log.Printf("Send error: %v\n", err)
	}
}

func (t *TgBot) Heic2JpgCompress(message *tgbotapi.Message) {
	file, err := t.GetDocument(message.Document.FileID)
	if err != nil {
		log.Printf("GetDocument error: %v\n", err)
		return
	}

	jpgFile, err := HeicToJpg(file)
	if err != nil {
		log.Printf("HeicToJpg error: %v\n", err)
		return
	}
	jpgName := strings.TrimRight(message.Document.FileName, "heic")
	jpgName = strings.TrimRight(jpgName, "HEIC")
	jpgName = fmt.Sprintf("%s.jpg", jpgName)

	imageToSend := tgbotapi.NewPhoto(message.Chat.ID, tgbotapi.FileBytes{
		Name:  jpgName,
		Bytes: jpgFile,
	})

	_, err = t.BotApi.Send(imageToSend)
	if err != nil {
		log.Printf("Send error: %v\n", err)
	}
}
