package main

import "net/http"

func main() {
	// 注册路由
	server := NewHttpServer("test-server", MetricsFilterBuilder)
	server.Route(http.MethodGet, "/", home)
	server.Route(http.MethodGet, "/body/once", readBodyOnce)
	server.Route(http.MethodGet, "/body/multi", getBodyIsNil)
	server.Route(http.MethodGet, "/wholeUrl", wholeUrl)
	server.Route(http.MethodGet, "/url/*", queryParams) // 支持通配符*匹配
	server.Route(http.MethodGet, "/header", header)
	server.Route(http.MethodGet, "/form", form)
	server.Route(http.MethodGet, "/user/signup", SignUp)

	err := server.Start(":8080")
	if err != nil {
		panic(err)
	}
}
