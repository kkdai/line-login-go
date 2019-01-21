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
var channelID, channelSecret string

//LINE MessageAPI related configuration
var serverURL string
var botToken, botSecret string

func main() {
	var err error
	serverURL = os.Getenv("LINECORP_PLATFORM_SERVERURL")
	channelID = os.Getenv("LINECORP_PLATFORM_CHANNEL_CHANNELID")
	channelSecret = os.Getenv("LINECORP_PLATFORM_CHANNEL_CHANNELSECRET")

	if bot, err = linebot.New(os.Getenv("LINECORP_PLATFORM_CHATBOT_CHANNELSECRET"), os.Getenv("LINECORP_PLATFORM_CHATBOT_CHANNELTOKEN")); err != nil {
		log.Println("Bot:", bot, " err:", err)
		return
	}

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", browse)

	//provide by Heroku
	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, nil)
}
