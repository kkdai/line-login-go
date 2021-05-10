package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	social "github.com/kkdai/line-login-sdk-go"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

var bot *linebot.Client

//LINE Login related configuration
var channelID, channelSecret string

//LINE MessageAPI related configuration
var serverURL string
var botToken, botSecret string
var socialClient *social.Client

func main() {
	var err error
	serverURL = os.Getenv("LINECORP_PLATFORM_SERVERURL")
	channelID = os.Getenv("LINECORP_PLATFORM_CHANNEL_CHANNELID")
	channelSecret = os.Getenv("LINECORP_PLATFORM_CHANNEL_CHANNELSECRET")

	if bot, err = linebot.New(os.Getenv("LINECORP_PLATFORM_CHATBOT_CHANNELSECRET"), os.Getenv("LINECORP_PLATFORM_CHATBOT_CHANNELTOKEN")); err != nil {
		log.Println("Bot:", bot, " err:", err)
		return
	}

	if socialClient, err = social.New(channelID, channelSecret); err != nil {
		log.Println("Social SDK:", socialClient, " err:", err)
		return
	}

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	//For LINE login
	http.HandleFunc("/", browse)
	http.HandleFunc("/gotoauthOpenIDpage", gotoauthOpenIDpage)
	http.HandleFunc("/gotoauthpage", gotoauthpage)
	http.HandleFunc("/auth", auth)

	//For linked chatbot
	http.HandleFunc("/callback", callbackHandler)

	//provide by Heroku
	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, nil)
}
