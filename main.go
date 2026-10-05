package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	social "github.com/kkdai/line-login-sdk-go"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

var bot *linebot.Client

// LINE Login related configuration
var channelID, channelSecret string

// LINE MessageAPI related configuration
var serverURL string
var botToken, botSecret string
var socialClient *social.Client

func main() {
	var err error
	serverURL = strings.TrimRight(os.Getenv("LINECORP_PLATFORM_SERVERURL"), "/")
	channelID = os.Getenv("LINECORP_PLATFORM_CHANNEL_CHANNELID")
	channelSecret = os.Getenv("LINECORP_PLATFORM_CHANNEL_CHANNELSECRET")

	if bot, err = linebot.New(os.Getenv("LINECORP_PLATFORM_CHATBOT_CHANNELSECRET"), os.Getenv("LINECORP_PLATFORM_CHATBOT_CHANNELTOKEN")); err != nil {
		log.Fatal("Bot: ", err)
	}

	if socialClient, err = social.New(channelID, channelSecret); err != nil {
		log.Fatal("Social SDK: ", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	//For LINE login
	mux.HandleFunc("/", browse)
	mux.HandleFunc("/gotoauthOpenIDpage", gotoauthOpenIDpage)
	mux.HandleFunc("/gotoauthpage", gotoauthpage)
	mux.HandleFunc("/auth", auth)

	//For linked chatbot
	mux.HandleFunc("/callback", callbackHandler)

	//PORT is provided by Cloud Run / Heroku
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%s", port),
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       120 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
