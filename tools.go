package main

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	IDToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
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

func RequestLoginToken(code, redirectURL, clientID, clientSecret string) (*TokenResponse, error) {
	qURL := url.QueryEscape(redirectURL)
	body := strings.NewReader(fmt.Sprintf("grant_type=authorization_code&code=%s&redirect_uri=%s&client_id=%s&client_secret=%s", code, qURL, clientID, clientSecret))
	req, err := http.NewRequest("POST", "https://api.line.me/oauth2/v2.1/token", body)
	if err != nil {
		// handle err
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
		return nil, err
	}
	defer resp.Body.Close()

	retBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("er:", err)
		return nil, err
	}
	retToken := TokenResponse{}
	if err := json.Unmarshal(retBody, &retToken); err != nil {
		return nil, err
	}

	return &retToken, nil
}

func VerifyToken(code, redirectURL, clientID, clientSecret string) (*TokenResponse, error) {
	qURL := url.QueryEscape(redirectURL)
	body := strings.NewReader(fmt.Sprintf("grant_type=authorization_code&code=%s&redirect_uri=%s&client_id=%s&client_secret=%s", code, qURL, clientID, clientSecret))
	req, err := http.NewRequest("POST", "https://api.line.me/oauth2/v2.1/token", body)
	if err != nil {
		// handle err
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
		return nil, err
	}
	defer resp.Body.Close()

	retBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("er:", err)
		return nil, err
	}
	retToken := TokenResponse{}
	if err := json.Unmarshal(retBody, &retToken); err != nil {
		return nil, err
	}

	return &retToken, nil
}
