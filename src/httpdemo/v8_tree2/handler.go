package main

type Handler interface {
	ServeHTTP(c *Context)
	Routeable //负责路由
}
type handlerFunc func(c *Context)
