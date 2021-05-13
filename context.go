package yong

import (
	"net/http"
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
