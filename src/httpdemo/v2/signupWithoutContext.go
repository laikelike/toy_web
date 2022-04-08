package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// SignUp:在没有 context 抽象的情况下，是长这样的
func SignUpOld(w http.ResponseWriter, r *http.Request) {
	req := &signUpReq{}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "read body failed: %v", err)
		// 要返回掉，不然就会继续执行后面的代码
		return
	}
	// 或返回json串,不优雅，不美观，还没业务逻辑就这么多了，体现抽象 重要性
	err = json.Unmarshal(body, req)
	if err != nil {
		fmt.Fprintf(w, "deserialized failed: %v", err)
		// 要返回掉，不然就会继续执行后面的代码
		return
	}

	resp := &commonResponse{
		Data: 123,
	}
	// 封装成writeJson，还可以进一步封装
	//err = ctx.WriteJson(http.StatusOK, resp)
	respJson, err := json.Marshal(resp)
	//如果写入失败，写err也很有可能会失败，无法给用户返回错误信息，只能自己记录日志
	//ctx.W.Write([]byte(err.Error()))
	if err != nil {
		fmt.Printf("写入响应失败：%v", err) // 只能记录日志
	}
	// 返回一个虚拟的 user id ,123表示注册成功了
	fmt.Fprintf(w, "%s", string(respJson))
}

/*
大量的干扰代码
例如：
传入参数可以聚合在一起：w http.ResponseWriter, r *http.Request--》ctx context

*/
