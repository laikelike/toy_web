package main

import (
	"net/http"
)

// Routeable 可路由的接口，让server和Handler都组合此接口
type Routeable interface {
	Route(method string, pattern string, handleFunc handlerFunc) error
}
type Server interface {
	Routeable
	Start(address string) error
}

// sdkHttpServer 基于net/http库实现
type sdkHttpServer struct {
	// Name server 的名字，给个标记，日志输出的时候用得上
	Name    string
	handler Handler
	root    Filter
}

// Route 注册路由，这个核心函数只依赖于一些很抽象的函数
func (s *sdkHttpServer) Route(method string, pattern string, handlerFunc handlerFunc) error{
	return s.handler.Route(method, pattern, handlerFunc)
}

func (s *sdkHttpServer) Start(address string) error {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		c := NewContext(writer, request)
		s.root(c)
	})
	return http.ListenAndServe(address, nil)
}

// NewHttpServer 返回实例，隐藏创建实例的细节
func NewHttpServer(name string, builders ...FilterBuilder) Server { // ...不定参数
	// 改用自己写的树
	handler := NewHandlerBasedOnTree()
	//handler := NewHandlerBaseOnMap()
	var root Filter = handler.ServeHTTP

	for i := len(builders) - 1; i >= 0; i-- {
		b := builders[i]
		root = b(root)
	}

	return &sdkHttpServer{
		Name:    name,
		handler: handler,
		root:    root,
	}
}
