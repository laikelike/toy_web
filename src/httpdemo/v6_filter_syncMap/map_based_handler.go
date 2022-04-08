package main

import (
	"net/http"
	"sync"
)

// Routeable 可路由的接口，让server和Handler都组合此接口
type Routeable interface {
	Route(method string, pattern string, handlerFunc func(ctx *Context))
}

type Handler interface {
	ServeHTTP(c *Context)
	Routeable //负责路由
}

type HandlerBasedOnMap struct {
	handlers sync.Map // key = method + url
}

// Route 具体实现
func (h *HandlerBasedOnMap) Route(method string, pattern string, handlerFunc func(ctx *Context)) {
	key := h.key(method, pattern)
	h.handlers.Store(key, handlerFunc)
}

// ServeHTTP 实现的是系统内的http包内的Handler接口的方法ServeHttp
func (h *HandlerBasedOnMap) ServeHTTP(c *Context) {
	key := h.key(c.R.Method, c.R.URL.Path)
	if handler, ok := h.handlers.Load(key); ok { //检查是否注册过
		handler.(func(c *Context))(c)
	} else { // 没找到
		c.W.WriteHeader(http.StatusNotFound)
		c.W.Write([]byte("Not Found"))
	}
}

func (h *HandlerBasedOnMap) key(method string, pattern string) string {
	return method + "#" + pattern
}

func NewHandlerBaseOnMap() Handler {
	return &HandlerBasedOnMap{
		handlers: sync.Map{},
	}
}

// 确保HandlerBasedOnMap一定实现了Handler接口,如果Handler接口发生变更就可以及时发现
var _ Handler = &HandlerBasedOnMap{}

/*
PUT /user 创建用户
POST /user 更新用户
DELETE /user 删除用户
GET /user 获取用户

Restful风格
决定动作       决定资源
http method + http path = http handler
*/
