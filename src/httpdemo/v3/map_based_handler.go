package main

import (
	"fmt"
	"net/http"
)

/*
PUT /user 创建用户
POST /user 更新用户
DELETE /user 删除用户
GET /user 获取用户

Restful风格
决定动作       决定资源
http method + http path = http handler
*/

type HandlerBasedOnMap struct {
	handlers map[string]func(ctx *Context) // key = method + url
}

// ServeHTTP 实现的是系统内的http包内的Handler接口的方法ServeHttp
func (h *HandlerBasedOnMap) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	key := h.key(request.Method, request.URL.Path)
	if handler, ok := h.handlers[key]; ok { //检查是否注册过
		handler(NewContext(writer, request))
	} else { // 没找到
		writer.WriteHeader(http.StatusNotFound) //返回404
		writer.Write([]byte("Not Found"))
	}
}

func (h *HandlerBasedOnMap) key(method string, pattern string) string {
	return fmt.Sprintf("%s#%s", method, pattern)
}

/*
基于map的路由有什么缺陷：
Route 方法依赖于知道HandlerBasedOnMap的内部细节，当我们需要利用路由树来实现的
时候，sdkHttpServer也要修改
*/
