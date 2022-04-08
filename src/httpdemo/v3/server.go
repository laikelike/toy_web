package main

import (
	"net/http"
)

type Server interface {
	// method:POST GET PUT之类的,用map的键装method达到约束方法的效果
	Route(method string, pattern string, handleFunc func(ctx *Context))
	Start(address string) error
}

// sdkHttpServer 基于net/http库实现
type sdkHttpServer struct {
	Name    string // Name server 的名字，给个标记，日志输出的时候用得上
	handler *HandlerBasedOnMap
}

// Route 注册路由，这个核心函数只依赖于一些很抽象的函数
func (s *sdkHttpServer) Route(method string, pattern string, handleFunc func(ctx *Context)) {
	key := s.handler.key(method, pattern)
	s.handler.handlers[key] = handleFunc
	// a.b.c这种情况就需要知道b的实现是map才能如此使用，建议只有一层如a.b
}

func (s *sdkHttpServer) Start(address string) error {
	http.Handle("/", s.handler) //只用初始注册一次,s.handler负责路由
	return http.ListenAndServe(address, nil)
}

// NewHttpServer 返回实例，隐藏创建实例的细节
func NewHttpServer(name string) Server {
	return &sdkHttpServer{
		Name: name,
	}
}
