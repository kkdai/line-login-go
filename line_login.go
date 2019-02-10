package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
)

var nounce string
var state string

const lineLoginURL string = "https://access.line.me/oauth2/v2.1/authorize?response_type=code"

func browse(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("login.tmpl"))
	if err := tmpl.Execute(w, nil); err != nil {
		log.Println("Template err:", err)
	}
}

func gotoauthpage(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Printf("ParseForm() err: %v\n", err)
		return
	}
	chatbot := r.FormValue("chatbot")

	scope := "profile openid" //profile | openid | email
	state = GenerateNounce()
	nounce = GenerateNounce()
	redirectURL := fmt.Sprintf("%s/auth", serverURL)
	targetURL := socialClient.GetWebLoinURL(redirectURL, state, scope, nounce, chatbot)
	http.Redirect(w, r, targetURL, http.StatusSeeOther)
}

func auth(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Printf("ParseForm() err: %v\n", err)
		return
	}
	code := r.FormValue("code")
	inState := r.FormValue("state")
	//Check the state
	if strings.Compare(state, inState) != 0 {
		log.Println("State is not matching.")
		return
	}
	friendshipStatusChanged := r.FormValue("friendship_status_changed")
	log.Println("code:", code, " state:", state, "friend status:", friendshipStatusChanged)

	//Request for access token
	IDToken, err := socialClient.GetAccessToken(fmt.Sprintf("%s/auth", serverURL), code).Do()
	if err != nil {
		log.Println("RequestLoginToken err:", err)
		return
	}

	//Decode IDToken
	payload, err := DecodeIDToken(IDToken.IDToken, channelID)
	if err != nil {
		log.Println("DecodeIDToken err:", err)
		return
	}

	//verify access token
	tmpl := template.Must(template.ParseFiles("login_success.tmpl"))
	if err := tmpl.Execute(w, payload); err != nil {
		log.Println("Template err:", err)
	}
}
