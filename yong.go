package yong

import (
	"fmt"
	"net/http"
	"os"
	"reflect"
	"runtime"
)

type funcMap map[string]HandlerFunc

type pathMap map[string]funcMap

type mWMap map[string][]HandlerFunc
type Router struct {
	method
	path pathMap
	middleWare
	mWMap mWMap
}

func Default() *Router {
	router := &Router{}
	router.path = make(pathMap)
	router.mWMap = make(mWMap)
	return router
}

type HandlerFunc func(*Context)

type method interface {
	GET(string, HandlerFunc)
	POST(string, HandlerFunc)
	PATCH(string, HandlerFunc)
	PUT(string, HandlerFunc)
	DELETE(string, HandlerFunc)
}

type middleWare interface {
	USE(string, ...HandlerFunc)
}

func (rt Router) USE(path string, handlers ...HandlerFunc) {
	if rt.mWMap[path] == nil {
		rt.mWMap[path] = handlers
	} else {
		for _, value := range rt.mWMap[path] {
			rt.mWMap[path] = append(rt.mWMap[path], value)
		}
	}
}

func (fMap funcMap) setMethod(method string, handler HandlerFunc) bool {
	if fMap[method] == nil {
		fMap[method] = handler
		return true
	}
	// 덮어쓰기 가능하게 하려면 여기서 처리해줘야함
	// 현재는 같은 path에 중복된 method타입도 허용 안됨.
	fmt.Println("같은 경로의 같은 요청타입이 존재합니다.")
	return false
}

func (pMap pathMap) setPath(path string) {
	if pMap[path] == nil {
		pMap[path] = make(funcMap)
	}
}

func (rt Router) setHandle(path string, method string, handler HandlerFunc) {
	rt.path.setPath(path)
	status := rt.path[path].setMethod(method, handler)
	if status {
		nuMiddleware := len(rt.mWMap[path])
		handlerName := runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()
		debugPrint("%-6s %-10s --> %s (%d handlers)\n", method, path, handlerName, nuMiddleware+1)
	}
}

func (rt *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	handler := rt.path[path][r.Method]
	middlewares := rt.mWMap[path]

	var c Context

	defer r.Body.Close()

	c.Writer = w
	c.Request = r
	for _, mHandler := range middlewares {
		mHandler(&c)
	}
	if handler == nil {
		if r.Method == "OPTIONS" {
			return
		} else {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 page not found")
			return
		}
	}
	handler(&c)
}

func (rt Router) GET(relativePath string, handler HandlerFunc) {
	rt.setHandle(relativePath, http.MethodGet, handler)
}
func (rt Router) POST(relativePath string, handler HandlerFunc) {
	rt.setHandle(relativePath, http.MethodPost, handler)
}
func (rt Router) PATCH(relativePath string, handler HandlerFunc) {
	rt.setHandle(relativePath, http.MethodPatch, handler)
}
func (rt Router) PUT(relativePath string, handler HandlerFunc) {
	rt.setHandle(relativePath, http.MethodPut, handler)
}
func (rt Router) DELETE(relativePath string, handler HandlerFunc) {
	rt.setHandle(relativePath, http.MethodDelete, handler)
}
func (rt Router) OPTIONS(relativePath string, handler HandlerFunc) {
	rt.setHandle(relativePath, http.MethodOptions, handler)
}

func resolveAddress(addr []string) string {
	switch len(addr) {
	case 0:
		if port := os.Getenv("PORT"); port != "" {
			return ":" + port
		}
		return ":8080"
	case 1:
		return addr[0]
	default:
		panic("too many parameters")
	}
}

func (rt *Router) Run(addr ...string) (err error) {
	address := resolveAddress(addr)
	fmt.Printf("Listening and serving HTTP on %s\n", address)
	err = http.ListenAndServe(address, rt)
	return
}
