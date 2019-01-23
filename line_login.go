package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	b64 "encoding/base64"

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

	payload, sig, err := jwt.Parse(IDToken.IDToken)
	if err != nil {
		log.Println("jwt.Parse err:", err, payload, sig)
		return

	}

	log.Println("jwt.Parse succeess:", string(payload), string(sig))

	rem := len(payload) % 4
	payloadStr := string(payload)
	log.Println("side of payload=", rem)
	if rem > 0 {
		i := 4 - rem
		for ; i >= 0; i-- {
			payloadStr = payloadStr + "="
		}
	}

	bPayload, err := b64.StdEncoding.DecodeString(payloadStr)
	if err != nil {
		log.Println("base64 decode err:", err)
		return
	}
	log.Println("base64 decode succeess:", string(bPayload))

	//verify access token
	tmpl := template.Must(template.ParseFiles("login_success.tmpl"))
	if err := tmpl.Execute(w, IDToken); err != nil {
		log.Println("Template err:", err)
	}
}
