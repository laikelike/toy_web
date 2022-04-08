package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

/*
上下文依赖于框架自身创建
*/
type Context struct { // 判断创建接口还是结构体，看有没有别的实现，没有的就结构体
	W http.ResponseWriter
	R *http.Request
}

// ReadJson:读出body并反序列化
func (c *Context) ReadJson(req interface{}) error { // interface{}：空结构体可以接收任何类型参数
	r := c.R
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(c.W, "read body failed: %v", err)
		// 要返回掉，不然就会继续执行后面的代码
		return err
	}
	err = json.Unmarshal(body, req) //将body中的字节反序列化成req格式
	if err != nil {
		return err
	}
	return nil
}

// WriteJson 写入状态码
func (c *Context) WriteJson(statusCode int, resp interface{}) error {

	c.W.WriteHeader(statusCode)
	respJson, err := json.Marshal(resp)
	if err != nil {
		fmt.Fprintf(c.W, "Resp Marshal failed: %v", err)
		return err
	}
	_, err = c.W.Write(respJson) // 第一个返回值代表写入多少数据
	return err

}

func (c *Context) OkJson(resp interface{}) error {
	return c.WriteJson(http.StatusOK, resp) // 返回200
}
func (c *Context) SystemErrorJson(resp interface{}) error {
	return c.WriteJson(http.StatusInternalServerError, resp) // 返回500
}
func (c *Context) BadRequsetJson(data interface{}) error {
	return c.WriteJson(http.StatusBadRequest, data) // 返回400
}

// NewContext不希望server理解创建ctx的细节，进一步对创建context封装
func NewContext(writer http.ResponseWriter, request *http.Request) *Context {
	return &Context{
		W: writer,
		R: request,
	}
}
