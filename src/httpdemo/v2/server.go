package main

import (
	"net/http"
)

/*
Server包表达一种逻辑上的抽象，它代表的是对某个端口进行监听的实体，
必要时可以开启多个Server，监听多个端口
*/
type Server interface {
	Route(pattern string, HandleFunc func(ctx *Context))
	Start(address string) error
}

// sdkHttpServer 基于net/http库实现
type sdkHttpServer struct {
	Name string // 标记不同Server
}

type Header map[string][]string

// Route 注册路由，命中该路由的会执行handlerFunc代码
func (s *sdkHttpServer) Route(pattern string, handleFunc func(ctx *Context)) {
	http.HandleFunc(pattern, func(writer http.ResponseWriter, request *http.Request) {
		ctx := NewContext(writer, request)
		handleFunc(ctx)
	})
}

// Start 启动服务器
func (s *sdkHttpServer) Start(address string) error {
	return http.ListenAndServe(address, nil)
}

// NewHttpServer 返回实例，隐藏创建实例的细节
func NewHttpServer(name string) Server {
	return &sdkHttpServer{
		Name: name,
	}
}
