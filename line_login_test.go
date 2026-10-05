package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	social "github.com/kkdai/line-login-sdk-go"
)

func setup(t *testing.T) {
	t.Helper()
	var err error
	serverURL = "https://example.com"
	channelID = "123"
	if socialClient, err = social.New(channelID, "secret"); err != nil {
		t.Fatal(err)
	}
}

func cookies(rec *httptest.ResponseRecorder) map[string]*http.Cookie {
	m := map[string]*http.Cookie{}
	for _, c := range rec.Result().Cookies() {
		m[c.Name] = c
	}
	return m
}

func TestStartLoginSetsCookiesAndRedirects(t *testing.T) {
	setup(t)
	rec := httptest.NewRecorder()
	gotoauthOpenIDpage(rec, httptest.NewRequest("GET", "/gotoauthOpenIDpage?chatbot=normal", nil))

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("status = %d", rec.Code)
	}
	loc, err := url.Parse(rec.Header().Get("Location"))
	if err != nil {
		t.Fatal(err)
	}
	cs := cookies(rec)
	if cs[cookieState] == nil || cs[cookieNonce] == nil {
		t.Fatalf("missing cookies: %v", cs)
	}
	if loc.Query().Get("state") != cs[cookieState].Value || loc.Query().Get("nonce") != cs[cookieNonce].Value {
		t.Error("redirect state/nonce do not match cookies")
	}
	if loc.Query().Get("bot_prompt") != "normal" {
		t.Errorf("bot_prompt = %q", loc.Query().Get("bot_prompt"))
	}
	if !cs[cookieState].HttpOnly || !cs[cookieState].Secure || cs[cookieState].Path != cookiePath {
		t.Error("state cookie not HttpOnly/Secure/scoped")
	}
	if cs[cookieVerifier] != nil || loc.Query().Get("code_challenge") != "" {
		t.Error("PKCE data present without enablePKCE")
	}
}

func TestStartLoginStateIsPerRequest(t *testing.T) {
	setup(t)
	var states []string
	for i := 0; i < 2; i++ {
		rec := httptest.NewRecorder()
		gotoauthpage(rec, httptest.NewRequest("GET", "/gotoauthpage", nil))
		states = append(states, cookies(rec)[cookieState].Value)
	}
	if states[0] == states[1] {
		t.Error("state reused across logins")
	}
}

func TestStartLoginPKCE(t *testing.T) {
	setup(t)
	rec := httptest.NewRecorder()
	gotoauthpage(rec, httptest.NewRequest("GET", "/gotoauthpage?enablePKCE=true", nil))
	loc, _ := url.Parse(rec.Header().Get("Location"))
	v := cookies(rec)[cookieVerifier]
	if v == nil {
		t.Fatal("missing verifier cookie")
	}
	if loc.Query().Get("code_challenge") != social.PkceChallenge(v.Value) {
		t.Error("code_challenge does not match verifier")
	}
}

func TestStartLoginIgnoresInvalidChatbot(t *testing.T) {
	setup(t)
	rec := httptest.NewRecorder()
	gotoauthpage(rec, httptest.NewRequest("GET", "/gotoauthpage?chatbot=evil", nil))
	loc, _ := url.Parse(rec.Header().Get("Location"))
	if loc.Query().Get("bot_prompt") != "" {
		t.Error("invalid chatbot value was forwarded")
	}
}

func TestAuthRejectsBadState(t *testing.T) {
	setup(t)
	cases := map[string]struct {
		query  string
		cookie string
	}{
		"empty both":     {"/auth?state=&code=x", ""},
		"no cookie":      {"/auth?state=abc&code=x", ""},
		"no query state": {"/auth?code=x", "abc"},
		"mismatch":       {"/auth?state=abc&code=x", "def"},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tc.query, nil)
			if tc.cookie != "" {
				req.AddCookie(&http.Cookie{Name: cookieState, Value: tc.cookie})
			}
			rec := httptest.NewRecorder()
			auth(rec, req)
			if rec.Code != http.StatusBadRequest {
				t.Errorf("status = %d, want 400", rec.Code)
			}
		})
	}
}

func TestAuthRejectsAuthorizationError(t *testing.T) {
	setup(t)
	req := httptest.NewRequest("GET", "/auth?state=abc&error=access_denied", nil)
	req.AddCookie(&http.Cookie{Name: cookieState, Value: "abc"})
	rec := httptest.NewRecorder()
	auth(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rec.Code)
	}
}

func TestAuthRejectsMissingCode(t *testing.T) {
	setup(t)
	req := httptest.NewRequest("GET", "/auth?state=abc", nil)
	req.AddCookie(&http.Cookie{Name: cookieState, Value: "abc"})
	rec := httptest.NewRecorder()
	auth(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rec.Code)
	}
}

func TestAuthClearsCookies(t *testing.T) {
	setup(t)
	rec := httptest.NewRecorder()
	auth(rec, httptest.NewRequest("GET", "/auth", nil))
	if c := cookies(rec)[cookieState]; c == nil || c.MaxAge >= 0 {
		t.Error("state cookie not cleared")
	}
}

func TestBrowseOnlyServesRoot(t *testing.T) {
	rec := httptest.NewRecorder()
	browse(rec, httptest.NewRequest("GET", "/nope", nil))
	if rec.Code != http.StatusNotFound {
		t.Errorf("status = %d, want 404", rec.Code)
	}
	rec = httptest.NewRecorder()
	browse(rec, httptest.NewRequest("GET", "/", nil))
	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", rec.Code)
	}
}
