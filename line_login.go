package main

import (
	"crypto/subtle"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	social "github.com/kkdai/line-login-sdk-go"
)

// Templates are parsed once at startup so a missing file fails fast.
var (
	loginTmpl   = template.Must(template.ParseFiles("login.tmpl"))
	successTmpl = template.Must(template.ParseFiles("login_success.tmpl"))
)

// Per-login values live in short-lived cookies scoped to /auth, so concurrent
// users never share state and the state is bound to the browser that started it.
const (
	cookieState    = "line_login_state"
	cookieNonce    = "line_login_nonce"
	cookieVerifier = "line_login_verifier"
	cookiePath     = "/auth"
	cookieMaxAge   = 600
)

func setLoginCookie(w http.ResponseWriter, name, value string) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     cookiePath,
		MaxAge:   cookieMaxAge,
		HttpOnly: true,
		Secure:   strings.HasPrefix(serverURL, "https://"),
		SameSite: http.SameSiteLaxMode,
	})
}

func clearLoginCookie(w http.ResponseWriter, name string) {
	http.SetCookie(w, &http.Cookie{Name: name, Value: "", Path: cookiePath, MaxAge: -1})
}

func cookieValue(r *http.Request, name string) string {
	c, err := r.Cookie(name)
	if err != nil {
		return ""
	}
	return c.Value
}

func browse(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	if err := loginTmpl.Execute(w, nil); err != nil {
		log.Println("Template err:", err)
	}
}

func gotoauthpage(w http.ResponseWriter, r *http.Request) {
	startLogin(w, r, "profile")
}

func gotoauthOpenIDpage(w http.ResponseWriter, r *http.Request) {
	startLogin(w, r, "profile openid")
}

// startLogin generates state/nonce (and a PKCE verifier on request), stores them
// in cookies, and redirects the user to the LINE authorization page.
func startLogin(w http.ResponseWriter, r *http.Request, scope string) {
	if err := r.ParseForm(); err != nil {
		log.Printf("ParseForm() err: %v\n", err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	// bot_prompt only accepts these values; ignore anything else.
	chatbot := r.FormValue("chatbot")
	if chatbot != "normal" && chatbot != "aggressive" {
		chatbot = ""
	}

	state, err := social.GenerateNonce()
	if err != nil {
		log.Println("GenerateNonce err:", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	nonce, err := social.GenerateNonce()
	if err != nil {
		log.Println("GenerateNonce err:", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	redirectURL := fmt.Sprintf("%s/auth", serverURL)
	opts := social.AuthRequestOptions{Nonce: nonce, BotPrompt: chatbot, Prompt: "consent"}

	var targetURL string
	if r.FormValue("enablePKCE") == "true" {
		verifier, err := social.GenerateCodeVerifier(64)
		if err != nil {
			log.Println("GenerateCodeVerifier err:", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		setLoginCookie(w, cookieVerifier, verifier)
		targetURL, err = socialClient.GetPKCEWebLoginURL(redirectURL, state, scope, social.PkceChallenge(verifier), opts)
		if err != nil {
			log.Println("GetPKCEWebLoginURL err:", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
	} else {
		targetURL, err = socialClient.GetWebLoginURL(redirectURL, state, scope, opts)
		if err != nil {
			log.Println("GetWebLoginURL err:", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
	}

	setLoginCookie(w, cookieState, state)
	setLoginCookie(w, cookieNonce, nonce)
	http.Redirect(w, r, targetURL, http.StatusSeeOther)
}

func auth(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Printf("ParseForm() err: %v\n", err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	// The cookies are single use: clear them whatever the outcome.
	state := cookieValue(r, cookieState)
	nonce := cookieValue(r, cookieNonce)
	verifier := cookieValue(r, cookieVerifier)
	clearLoginCookie(w, cookieState)
	clearLoginCookie(w, cookieNonce)
	clearLoginCookie(w, cookieVerifier)

	// Check the state: both sides must be present and equal.
	inState := r.FormValue("state")
	if state == "" || inState == "" || subtle.ConstantTimeCompare([]byte(state), []byte(inState)) != 1 {
		log.Println("State is not matching.")
		http.Error(w, "invalid state", http.StatusBadRequest)
		return
	}

	if e := r.FormValue("error"); e != "" {
		log.Println("Authorization error:", e)
		http.Error(w, "login was not authorized", http.StatusBadRequest)
		return
	}
	code := r.FormValue("code")
	if code == "" {
		http.Error(w, "missing code", http.StatusBadRequest)
		return
	}
	log.Println("friend status:", r.FormValue("friendship_status_changed"))

	//Request for access token
	redirectURL := fmt.Sprintf("%s/auth", serverURL)
	var token *social.TokenResponse
	var err error
	if verifier != "" {
		token, err = socialClient.GetAccessTokenPKCE(redirectURL, code, verifier).WithContext(r.Context()).Do()
	} else {
		token, err = socialClient.GetAccessToken(redirectURL, code).WithContext(r.Context()).Do()
	}
	if err != nil {
		log.Println("RequestLoginToken err:", err)
		http.Error(w, "failed to get access token", http.StatusBadGateway)
		return
	}

	//Verify the access token
	if _, err := socialClient.TokenVerify(token.AccessToken).WithContext(r.Context()).Do(); err != nil {
		log.Println("TokenVerify err:", err)
		http.Error(w, "failed to verify access token", http.StatusBadGateway)
		return
	}

	var payload *social.BasicPayload
	if len(token.IDToken) == 0 {
		// User didn't request openID, use access token to get user profile
		res, err := socialClient.GetUserProfile(token.AccessToken).WithContext(r.Context()).Do()
		if err != nil {
			log.Println("GetUserProfile err:", err)
			http.Error(w, "failed to get user profile", http.StatusBadGateway)
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
			http.Error(w, "invalid id token", http.StatusUnauthorized)
			return
		}
	}

	if err := successTmpl.Execute(w, payload); err != nil {
		log.Println("Template err:", err)
	}
}
