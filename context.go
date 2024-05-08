package yong

import (
	"net/http"
	"net/url"
	"sync"
)

type Context struct {
	Request *http.Request
	Writer  http.ResponseWriter
	mu      sync.RWMutex
	Keys    map[string]interface{}
}

func (c *Context) Set(key string, value interface{}) {
	c.mu.Lock()
	if c.Keys == nil {
		c.Keys = make(map[string]interface{})
	}

	c.Keys[key] = value
	c.mu.Unlock()
}

func (c *Context) Get(key string) (value interface{}, exists bool) {
	c.mu.RLock()
	value, exists = c.Keys[key]
	c.mu.RUnlock()
	return
}

func (c *Context) SetCookie(name, value string, maxAge int, path, domain string, SameSite http.SameSite, secure, httpOnly bool) {
	if path == "" {
		path = "/"
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     name,
		Value:    url.QueryEscape(value),
		MaxAge:   maxAge,
		Path:     path,
		Domain:   domain,
		SameSite: SameSite,
		Secure:   secure,
		HttpOnly: httpOnly,
	})
}
