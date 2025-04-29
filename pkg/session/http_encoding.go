package session

import (
	"encoding/json"
	"net/http"

	"github.com/euiko/webapp/pkg/helper"
)

type (
	HTTPCookieEncoding struct {
		path     string
		httpOnly bool
		secure   bool
		key      string
		r        *http.Request
		w        http.ResponseWriter
	}

	HTTPCookieEncodingOption func(*HTTPCookieEncoding)
)

func NewHTTPCookieEncoding(w http.ResponseWriter, r *http.Request, opts ...HTTPCookieEncodingOption) *HTTPCookieEncoding {
	e := HTTPCookieEncoding{
		path:     "/",
		httpOnly: false,
		secure:   false,
		key:      "session",
		r:        r,
		w:        w,
	}

	for _, opt := range opts {
		opt(&e)
	}

	return &e
}

func (e *HTTPCookieEncoding) Decode() (*Session, error) {
	var (
		cookie       *http.Cookie
		cookieHeader = e.r.Header.Get("Cookie")
		session      = New()
	)

	cookies, err := http.ParseCookie(cookieHeader)
	if err != nil {
		return session, err
	}

	for _, c := range cookies {
		if c.Name == e.key {
			cookie = c
			break
		}
	}
	// create new session if none cookie found
	if cookie == nil {
		return session, nil
	}

	jsoned, err := helper.DecodeBase64(cookie.Value)
	if err != nil {
		return session, err
	}

	values := make(map[string]any)
	err = json.Unmarshal(jsoned, &values)
	if err != nil {
		return session, err
	}

	for key, value := range values {
		session.Store(key, value)
	}

	return session, nil
}

func (e *HTTPCookieEncoding) Encode(session *Session) error {
	sessionMap := make(map[string]interface{})
	session.Range(func(key, value any) bool {
		sessionMap[key.(string)] = value
		return true
	})

	encoded, err := json.Marshal(sessionMap)
	if err != nil {
		return err
	}

	base64 := helper.EncodeBase64(encoded)
	cookie := http.Cookie{
		Name:     e.key,
		Value:    base64,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
	}

	http.SetCookie(e.w, &cookie)
	return nil
}
