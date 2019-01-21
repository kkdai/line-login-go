package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/line/line-bot-sdk-go/linebot"
)

var bot *linebot.Client

//LINE Login related configuration
var channelServerURL, channelID, channelSecret string

//LINE MessageAPI related configuration
var botServerURL string

func main() {
	var err error
	channelServerURL = os.Getenv("LINECORP_PLATFORM_CHANNEL_SERVERURL")
	channelID = os.Getenv("LINECORP_PLATFORM_CHANNEL_CHANNELID")
	channelSecret = os.Getenv("LINECORP_PLATFORM_CHANNEL_CHANNELSECRET")

	botServerURL = os.Getenv("LINECORP_PLATFORM_CHATBOT_SERVERURL")
	if bot, err = linebot.New(os.Getenv("LINECORP_PLATFORM_CHATBOT_CHANNELSECRET"), os.Getenv("LINECORP_PLATFORM_CHATBOT_CHANNELTOKEN")); err != nil {
		log.Println("Bot:", bot, " err:", err)
		return
	}

	//provide by Heroku
	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, nil)
}
