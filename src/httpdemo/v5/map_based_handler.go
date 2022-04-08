package main

import (
	"net/http"
)

// Routeable 可路由的接口，让server和Handler都组合此接口
type Routeable interface {
	Route(method string, pattern string, handlerFunc func(ctx *Context))
}

type Handler interface {
	http.Handler //这是系统内接口，负责处理请求
	Routeable    //负责路由
}

type HandlerBasedOnMap struct {
	handlers map[string]func(ctx *Context) // key = method + url
}

// Route 具体实现
func (h *HandlerBasedOnMap) Route(method string, pattern string, handlerFunc func(ctx *Context)) {
	key := h.key(method, pattern)
	h.handlers[key] = handlerFunc
}

// ServeHTTP 实现的是系统内的http包内的Handler接口的方法ServeHttp
func (h *HandlerBasedOnMap) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	key := h.key(request.Method, request.URL.Path)
	if handler, ok := h.handlers[key]; ok { //检查是否注册过
		handler(NewContext(writer, request))
	} else { // 没找到
		writer.WriteHeader(http.StatusNotFound)
		writer.Write([]byte("Not Found"))
	}
}

func (h *HandlerBasedOnMap) key(method string, pattern string) string {
	return method + "#" + pattern
}

func NewHandlerBaseOnMap() Handler {
	return &HandlerBasedOnMap{
		handlers: make(map[string]func(ctx *Context)),
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
