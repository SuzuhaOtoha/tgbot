package main

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gopkg.in/yaml.v2"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"tgbot/model"
	"time"
)

var conf = new(model.Config)

var StartTime time.Time
var Images []string
var SendTextTotal = 0
var SengImageTotal = 0

var Enabled = true
var Repeat = false

func login() *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(conf.Token)
	bot.Debug = conf.Debug
	if err != nil {
		panic(err)
	}
	return bot
}

//func loadImage() []string {
//	var files []string
//	err := filepath.Walk(conf.ImagePath, func(path string, info os.FileInfo, err error) error {
//		files = append(files, path)
//		return nil
//	})
//	if err != nil {
//		panic(err)
//	}
//	return files
//}

func sendText(bot *tgbotapi.BotAPI, chatID int64, text string) {
	if !Enabled {
		return
	}
	msg := tgbotapi.NewMessage(chatID, text)
	if _, err := bot.Send(msg); err != nil {
		log.Panic(err)
	}
	SendTextTotal++
	return
}

func sendPhoto(bot *tgbotapi.BotAPI, chatID int64, imagePath string) error {
	//if !Enabled {
	//	return nil
	//}
	//photo := tgbotapi.NewPhoto(chatID, tgbotapi.FilePath(imagePath))
	//if _, err := bot.Send(photo); err != nil {
	//	return nil
	//}
	//SengImageTotal++
	//return nil
	msg := tgbotapi.NewMessage(chatID, imagePath)
	msg.ParseMode = "Markdown"
	if _, err := bot.Send(msg); err != nil {
		log.Panic(err)
	}
	SengImageTotal++
	return nil
}

func getStatus() string {
	enable := ""
	if Enabled {
		enable = "true"
	} else {
		enable = "false"
	}
	repeat := ""
	if Repeat {
		repeat = "true"
	} else {
		repeat = "false"
	}
	text := fmt.Sprintf(
		`Enable:%s
Repeat:%s
Uptime:%s
Total number of Image:%d
Total number of send Text:%d
Total number of send Image:%d
Total number of send message:%d`,
		enable,
		repeat,
		time.Since(StartTime),
		len(Images),
		SendTextTotal,
		SengImageTotal,
		SengImageTotal+SendTextTotal,
	)
	return text
}

func isAdmin(chatID int64) bool {
	if chatID == conf.AdminAccount {
		return true
	}
	return false
}

func status(bot *tgbotapi.BotAPI, chatID int64) {
	if isAdmin(chatID) {
		msg := tgbotapi.NewMessage(chatID, getStatus())
		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
		SendTextTotal++
		return
	}
	sendText(bot, chatID, conf.Info)
	return

}

func setu(bot *tgbotapi.BotAPI, chatID int64) {
	//imageId := rand.Intn(len(Images) - 1)
	//filename := Images[imageId]
	var loliconApiRet model.LoliconApiRet
	for {
		resp, err := http.Get("https://api.lolicon.app/setu/v2?r18=2&tag=%E8%90%9D%E8%8E%89")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				return
			}
		}(resp.Body)
		body, _ := io.ReadAll(resp.Body)
		err = json.Unmarshal(body, &loliconApiRet)
		if err != nil {
			fmt.Println(err)
		}
		resp, err = http.Get(loliconApiRet.Data[0].Urls.Original)
		if resp.StatusCode == 200 {
			break
		}
	}
	title := loliconApiRet.Data[0].Title
	link := loliconApiRet.Data[0].Urls.Original
	mdString := fmt.Sprintf("[%s](%s)", title, link)
	err := sendPhoto(bot, chatID, mdString)
	if err != nil {
		setu(bot, chatID)
	}
	return
}

func getUpdate(bot *tgbotapi.BotAPI) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			continue
		}
		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "help":
				sendText(bot, update.Message.Chat.ID, conf.Info)
			case "status":
				status(bot, update.Message.Chat.ID)
			case "setu":
				setu(bot, update.Message.Chat.ID)
			case "start":
				if isAdmin(update.Message.Chat.ID) {
					Enabled = true
					sendText(bot, update.Message.Chat.ID, "bot start!")
					continue
				}
				sendText(bot, update.Message.Chat.ID, conf.Info)
			case "stop":
				if isAdmin(update.Message.Chat.ID) {
					sendText(bot, update.Message.Chat.ID, "bot stop!")
					Enabled = false
					continue
				}
				sendText(bot, update.Message.Chat.ID, conf.Info)
			case "repeaton":
				if isAdmin(update.Message.Chat.ID) {
					sendText(bot, update.Message.Chat.ID, "repeat on!")
					Repeat = true
					continue
				}
				sendText(bot, update.Message.Chat.ID, conf.Info)
			case "repeatoff":
				if isAdmin(update.Message.Chat.ID) {
					sendText(bot, update.Message.Chat.ID, "repeat off!")
					Repeat = false
					continue
				}
				sendText(bot, update.Message.Chat.ID, conf.Info)
			default:
				sendText(bot, update.Message.Chat.ID, conf.Info)
			}
		} else {
			if Repeat {
				sendText(bot, update.Message.Chat.ID, update.Message.Text)
				continue
			}
		}
	}
}

func loadConfig() {
	userPath, _ := os.UserHomeDir()
	configPath := ""
	sysType := runtime.GOOS
	if sysType == "windows" {
		configPath = userPath + "\\.config\\tgbot.yaml"
	}
	if sysType == "linux" {
		configPath = userPath + "/.config/tgbot.yaml"
	}
	log.Printf("Use config file %s\n", configPath)
	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(yamlFile, conf)
	if err != nil {
		panic(err)
	}
}

func main() {
	loadConfig()
	log.Printf("config load success!token=%s,image path=%s,admin account=%s", conf.Token, conf.ImagePath, strconv.FormatInt(conf.AdminAccount, 10))
	StartTime = time.Now()
	bot := login()
	log.Printf("Authorized on account %s", bot.Self.UserName)

	//Images = loadImage()
	//log.Printf("load %d image(s) from %s", len(Images), conf.ImagePath)

	getUpdate(bot)
}
