package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gbrlsnchs/jwt"
)

const lineLoginURL string = "https://access.line.me/oauth2/v2.1/authorize?response_type=code"

func browse(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("login.tmpl"))
	if err := tmpl.Execute(w, nil); err != nil {
		log.Println("Template err:", err)
	}
}

func gotoauthpage(w http.ResponseWriter, r *http.Request) {
	scope := "profile openid" //profile | openid | email
	state := GenerateNounce()
	clientID := channelID
	nounce := GenerateNounce()
	redirectURL := fmt.Sprintf("%s/auth", serverURL)
	targetURL := GetWebLoinURL(clientID, redirectURL, state, scope, nounce)
	log.Println("url=", targetURL)
	http.Redirect(w, r, targetURL, http.StatusSeeOther)
}

func auth(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Printf("ParseForm() err: %v\n", err)
		return
	}
	code := r.FormValue("code")
	state := r.FormValue("state")
	friendshipStatusChanged := r.FormValue("friendship_status_changed")
	log.Println("code:", code, " state:", state, "friend status:", friendshipStatusChanged)

	//Request for access token
	IDToken, err := RequestLoginToken(code, fmt.Sprintf("%s/auth", serverURL), channelID, channelSecret)
	if err != nil {
		log.Println("RequestLoginToken err:", err)
		return
	}

	byteIDToken, _ := json.Marshal(IDToken)
	log.Println("ID_Token:=", IDToken)
	payload, sig, err := jwt.Parse(string(byteIDToken))
	if err != nil {
		log.Println("jwt.Parse err:", err, payload, sig)
		return

	}

	log.Println("jwt.Parse succeess:", payload, sig)

	//verify access token
	tmpl := template.Must(template.ParseFiles("login_success.tmpl"))
	if err := tmpl.Execute(w, IDToken); err != nil {
		log.Println("Template err:", err)
	}
}
