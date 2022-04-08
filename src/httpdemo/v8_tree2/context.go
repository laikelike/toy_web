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
type Context struct {
	W http.ResponseWriter
	R *http.Request
}

// ReadJson 读出body并反序列化
func (c *Context) ReadJson(req interface{}) error {
	r := c.R
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(c.W, "read body failed: %v", err)
		// 要返回掉，不然就会继续执行后面的代码
		return err
	}

	err = json.Unmarshal(body, req)
	if err != nil {
		return err
	}
	return nil
}

func (c *Context) WriteJson(code int, resp interface{}) error {

	c.W.WriteHeader(code)
	respJson, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	_, err = c.W.Write(respJson)
	if err != nil {
		return err
	}
	return nil
}

func (c *Context) OkJson(resp interface{}) error {
	return c.WriteJson(http.StatusOK, resp)
}
func (c *Context) SystemErrorJson(resp interface{}) error {
	return c.WriteJson(http.StatusInternalServerError, resp)
}
func (c *Context) BadRequestJson(resp interface{}) error {
	return c.WriteJson(http.StatusBadRequest, resp)
}

// NewContext：不希望server理解创建ctx的细节，进一步对创建context封装
func NewContext(writer http.ResponseWriter, request *http.Request) *Context {
	return &Context{
		R: request,
		W: writer,
	}
}
