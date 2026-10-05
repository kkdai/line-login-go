package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	social "github.com/kkdai/line-login-sdk-go"
)

var nonce string
var state string

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

	scope := "profile" //profile | openid | email
	var err error
	if state, err = social.GenerateNonce(); err != nil {
		log.Println("GenerateNonce err:", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	if nonce, err = social.GenerateNonce(); err != nil {
		log.Println("GenerateNonce err:", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	redirectURL := fmt.Sprintf("%s/auth", serverURL)
	targetURL, err := socialClient.GetWebLoginURL(redirectURL, state, scope, social.AuthRequestOptions{Nonce: nonce, BotPrompt: chatbot, Prompt: "consent"})
	if err != nil {
		log.Println("GetWebLoginURL err:", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, targetURL, http.StatusSeeOther)
}

func gotoauthOpenIDpage(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Printf("ParseForm() err: %v\n", err)
		return
	}
	chatbot := r.FormValue("chatbot")

	scope := "profile openid" //profile | openid | email
	var err error
	if state, err = social.GenerateNonce(); err != nil {
		log.Println("GenerateNonce err:", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	if nonce, err = social.GenerateNonce(); err != nil {
		log.Println("GenerateNonce err:", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	redirectURL := fmt.Sprintf("%s/auth", serverURL)
	targetURL, err := socialClient.GetWebLoginURL(redirectURL, state, scope, social.AuthRequestOptions{Nonce: nonce, BotPrompt: chatbot, Prompt: "consent"})
	if err != nil {
		log.Println("GetWebLoginURL err:", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
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
	token, err := socialClient.GetAccessToken(fmt.Sprintf("%s/auth", serverURL), code).Do()
	if err != nil {
		log.Println("RequestLoginToken err:", err)
		return
	}

	log.Println("access_token:", token.AccessToken, " refresh_token:", token.RefreshToken)

	//Start to verify token and renew it.
	if result, err := socialClient.TokenVerify(token.AccessToken).Do(); err != nil {
		log.Println("TokenVerify err:", err, result)
		return
	}

	//Start to refresh token and renew it.
	if refresh, err := socialClient.RefreshToken(token.RefreshToken).Do(); err != nil {
		log.Println("RefreshToken err:", err, refresh)
		return
	}

	var payload *social.BasicPayload
	if len(token.IDToken) == 0 {
		// User don't request openID, use access token to get usere profile
		log.Println(" token:", token, " AccessToken:", token.AccessToken)
		res, err := socialClient.GetUserProfile(token.AccessToken).Do()
		if err != nil {
			log.Println("GetUserProfile err:", err)
			return
		}
		payload = &social.BasicPayload{
			Name:    res.DisplayName,
			Picture: res.PictureURL,
		}
	} else {
		//Decode token.IDToken to payload
		payload, err = token.DecodePayloadWithOptions(channelID, social.DecodePayloadOptions{Nonce: nonce, CheckExpiry: true})
		if err != nil {
			log.Println("DecodeIDToken err:", err)
			return
		}
	}

	//verify access token
	tmpl := template.Must(template.ParseFiles("login_success.tmpl"))
	if err := tmpl.Execute(w, payload); err != nil {
		log.Println("Template err:", err)
	}
}
