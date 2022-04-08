package main

import (
	"encoding/json"
	"fmt"
	"io"
)

func home(ctx *Context) {
	// Fprintf根据格式指定符进行格式化并写入w
	fmt.Fprintf(ctx.W, "Hi there, I love %s!", ctx.R.URL.Path[1:])
}

// readBodyOnce Body只能读一次,类似流的概念
func readBodyOnce(ctx *Context) {
	body, err := io.ReadAll(ctx.R.Body)
	if err != nil {
		fmt.Fprintf(ctx.W, "read body failed: %v", err)
		return
	}
	fmt.Fprintf(ctx.W, "read the data : %s \n", string(body))
	body, err = io.ReadAll(ctx.R.Body) // 尝试再次读取，啥也读不到，但是也不会报错
	if err != nil {                    // 不会进来这里
		fmt.Fprintf(ctx.W, "read the data one more time got error: %v", err)
		return
	}
	fmt.Fprintf(ctx.W, "read the data one more time: [%s] and read data length %d \n", string(body), len(body))
}
func getBodyIsNil(ctx *Context) {
	if ctx.R.GetBody == nil {
		fmt.Fprintf(ctx.W, "GetBody is nil \n")
	} else {
		fmt.Fprintf(ctx.W, "GetBody not nil \n")
	}
}

// queryParams：查询参数
func queryParams(ctx *Context) {
	values := ctx.R.URL.Query() // 返回是个map
	fmt.Fprintf(ctx.W, "query is %v!", values)
}

func wholeUrl(ctx *Context) {

	data, _ := json.Marshal(ctx.R.URL)
	fmt.Fprintf(ctx.W, "%s", string(data))
}

func header(ctx *Context) {
	fmt.Fprintf(ctx.W, "header is %v! \n ", ctx.R.Header)
}

func form(ctx *Context) {
	fmt.Fprintf(ctx.W, "before parse from %v!\n", ctx.R.Form)
	err := ctx.R.ParseForm()
	if err != nil {
		fmt.Fprintf(ctx.W, "parse form error %v\n", ctx.R.Form)
	}
	fmt.Fprintf(ctx.W, "after parse form %v \n", ctx.R.Form)
}

// signUpReq 传入json参数
type signUpReq struct {
	Email             string `json:"email"` // 运行时可以通过反射拿到email
	Password          string `json:"password"`
	ConfirmedPassword string `json:"confirmed_password"`
}

// commonResponse 返回json字符串
type commonResponse struct {
	BizCode int         `json:"biz_code"`
	Msg     string      `json:"msg"`
	Data    interface{} `json:"data"`
}

// SignUp 抽象优化，封装
func SignUp(ctx *Context) {
	req := &signUpReq{}

	err := ctx.ReadJson(req)
	if err != nil {
		ctx.BadRequsetJson(err)
		return
	}

	resp := &commonResponse{
		Data: 1234,
	}
	err = ctx.OkJson(resp)
	if err != nil {
		fmt.Printf("写入响应失败： %v", err)
	}
}

func main() {
	server := NewHttpServer("test-server")

	server.Route("/", home)
	server.Route("/body/once", readBodyOnce)
	server.Route("/body/multi", getBodyIsNil)
	server.Route("/url/query", queryParams)
	server.Route("/wholeUrl", wholeUrl)
	server.Route("/header", header)
	server.Route("/form", form)
	server.Route("/user/signup", SignUp)

	server.Start(":8080")

}
