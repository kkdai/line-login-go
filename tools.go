package main

import (
	b64 "encoding/base64"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

//GetWebLoinURL - LINE LOGIN 2.1 get LINE Login URL
func GetWebLoinURL(clientID, redirectURL, state, scope, nounce string) string {
	req, err := http.NewRequest("GET", "https://access.line.me/oauth2/v2.1/authorize", nil)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
	q := req.URL.Query()
	q.Add("response_type", "code")
	q.Add("client_id", clientID)
	q.Add("state", state)
	q.Add("scope", scope)
	q.Add("nounce", nounce)
	q.Add("redirect_uri", redirectURL)
	req.URL.RawQuery = q.Encode()
	log.Println(req.URL.String())
	return req.URL.String()
}

func GenerateNounce() string {
	return b64.StdEncoding.EncodeToString([]byte(RandStringRunes(8)))
}
