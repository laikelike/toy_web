package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func home(ctx *Context) {
	// Fprintf根据格式指定符进行格式化并写入w
	fmt.Fprintf(ctx.W, "Hi there, I love %s.", ctx.R.URL.Path[1:])
}

// readBodyOnce Body只能读一次,类似流的概念
func readBodyOnce(ctx *Context) {
	body, err := io.ReadAll(ctx.R.Body)
	if err != nil {
		fmt.Fprintf(ctx.W, "read body failed: %v", err)
		return
	}
	fmt.Fprintf(ctx.W, "read the data: %s \n", string(body))
	body, err = io.ReadAll(ctx.R.Body) // 尝试再次读取，啥也读不到，但是也不会报错
	if err != nil {                    // 不会进来这里
		fmt.Fprintf(ctx.W, "read the data one more time got error: %v", err)
		return
	}
	fmt.Fprintf(ctx.W, "read the data one more time: [%s] and read data length %d \n", string(body), len(body))
}
func getBodyIsNil(ctx *Context) {
	if ctx.R.GetBody == nil {
		fmt.Fprint(ctx.W, "GetBody is nil \n")
	} else {
		fmt.Fprintf(ctx.W, "GetBody not nil \n")
	}
}

// queryParams：查询参数
func wholeUrl(ctx *Context) {
	data, _ := json.Marshal(ctx.R.URL)
	fmt.Fprintf(ctx.W, "%s", string(data))
}

func queryParams(ctx *Context) {
	values := ctx.R.URL.Query()
	fmt.Fprintf(ctx.W, "query is %v\n", values)
}

func header(ctx *Context) {
	fmt.Fprintf(ctx.W, "header is %v \n", ctx.R.Header)
}

func form(ctx *Context) {
	fmt.Fprintf(ctx.W, "before parse form %v\n", ctx.R.Form)
	//使用表单前需要先调用parseForm
	err := ctx.R.ParseForm()
	if err != nil {
		fmt.Fprintf(ctx.W, "parse form error %v\n", ctx.R.Form)
	}
	fmt.Fprintf(ctx.W, "after parse form %v\n", ctx.R.Form)
}

type signUpReq struct {
	// `json:"email"`是Tag，运行时可以通过反射拿到，声明式写法
	Email             string `json:"email"`
	Password          string `json:"password"`
	ConfirmedPassword string `json:"confirmed_password"`
}

type commonResponse struct {
	BizCode int         `json:"biz_code"`
	Msg     string      `json:"msg"`
	Data    interface{} `json:"data"`
}

func SignUp(ctx *Context) {
	req := &signUpReq{}

	err := ctx.ReadJson(req)
	if err != nil {
		ctx.BadRequestJson(err)
		return
	}
	resp := &commonResponse{
		Data: 123,
	}
	err = ctx.WriteJson(http.StatusOK, resp)
	if err != nil {
		fmt.Printf("写入响应失败：%v", err)
	}
}

func main() {
	// 注册路由
	server := NewHttpServer("test-server")
	server.Route(http.MethodGet, "/", home)
	server.Route(http.MethodGet, "/body/once", readBodyOnce)
	server.Route(http.MethodGet, "/body/multi", getBodyIsNil)
	server.Route(http.MethodGet, "/wholeUrl", wholeUrl)
	server.Route(http.MethodGet, "/url/query", queryParams)
	server.Route(http.MethodGet, "/header", header)
	server.Route(http.MethodGet, "/form", form)
	server.Route(http.MethodGet, "/user/signup", SignUp)

	err := server.Start(":8080")
	if err != nil {
		panic(err)
	}
}
