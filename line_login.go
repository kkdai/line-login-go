package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

const lineLoginURL string = "https://access.line.me/oauth2/v2.1/authorize?response_type=code"

func browse(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("login.tmpl"))
	if err := tmpl.Execute(w, nil); err != nil {
		log.Println("Template err:", err)
	}
}

func gotoauthpage(w http.ResponseWriter, r *http.Request) {
	scope := "profile%20openid%20email" //profile | openid | email
	state := GenerateNounce()
	clientID := channelID
	nounce := GenerateNounce()
	redirectURL := fmt.Sprintf("%s/auth", serverURL)
	targetURL := GetWebLoinURL(clientID, redirectURL, state, scope, nounce)
	log.Println("url=", targetURL)
}
