package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	// Fprintf根据格式指定符进行格式化并写入w
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

// Body只能读一次
func readBodyOnce(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "read body failed: %v", err)
		return
	}
	fmt.Fprintf(w, "read the data : %s \n", string(body))
	body, err = io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "read the data one more time got error: %v", err)
		return
	}
	fmt.Fprintf(w, "read the data one more time: [%s] and read data length %d \n", string(body), len(body))
}
func getBodyIsNil(w http.ResponseWriter, r *http.Request) {
	// body, _ := r.GetBody()
	// io.ReadAll(body) // 读出nil

	// body, _ = r.GetBody()
	// io.ReadAll(body) // 读出来nil
	// GetBody：原则上可以多次读取，但是在原生的http.Request里面，这个是nil
	if r.GetBody == nil {
		fmt.Fprintf(w, "GetBody is nil \n")
	} else {
		fmt.Fprintf(w, "GetBody not nil \n")
	}
}

// query：查询参数
func queryParams(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query() // 返回是个map
	fmt.Fprintf(w, "query is %v!", values)
}

func wholeUrl(w http.ResponseWriter, r *http.Request) {
	/*
	   {
	       "Scheme": "",       // http或者https
	       "Opaque": "",
	       "User": null,
	       "Host": "",
	       "Path": "/wholeUrl",
	       "RawPath": "",
	       "ForceQuery": false,
	       "RawQuery": "",     // 参数"name=cuicui"
	       "Fragment": "",
	       "RawFragment": ""
	   }
	*/
	data, _ := json.Marshal(r.URL)
	fmt.Fprintf(w, "%s", string(data))
}

func header(w http.ResponseWriter, r *http.Request) {
	/*
	   header大体上两类，一类是http预定义的；一类是自己定义的
	   Go会自动将header名字转为标准名字(大小写调整)
	   一般X开头表明是自己定义的比如 X-mycompany-age=18
	*/

	// Fprintf根据格式指定符进行格式化并写入w
	fmt.Fprintf(w, "header is %v! \n ", r.Header)
}

func form(w http.ResponseWriter, r *http.Request) {
	// Fprintf根据格式指定符进行格式化并写入w
	fmt.Fprintf(w, "before parse from %v!\n", r.Form)
	err := r.ParseForm() // 使用form表单前调用parseform
	if err != nil {
		fmt.Fprintf(w, "parse form error %v\n", r.Form)
	}
	fmt.Fprintf(w, "after parse form %v \n", r.Form)
}

type signUpReq struct {
	Email             string `json:"email"` // 运行时可以通过反射拿到email
	Password          string `json:"password"`
	ConfirmedPassword string `json:"confirmed_password"`
}

type commonResponse struct {
	BizCode int         `json:"biz_code"`
	Msg     string      `json:"msg"`
	Data    interface{} `json:"data"`
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	req := &signUpReq{}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "read body failed: %v", err)
		// 要返回掉，不然就会继续执行后面的代码
		return
	}
	err = json.Unmarshal(body, req)
	if err != nil {
		fmt.Fprintf(w, "deserialized failed: %v", err)
		// 要返回掉，不然就会继续执行后面的代码
		return
	}

	resp := &commonResponse{
		Data: 123,
	}
	respJson, err := json.Marshal(resp)
	if err != nil {
		fmt.Fprintf(w, "Resp Marshal failed: %v", err)
		return
	}
	// 返回一个虚拟的 user id ,123表示注册成功了
	fmt.Fprintf(w, "%s", string(respJson))
}

func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/body/once", readBodyOnce)
	http.HandleFunc("/body/multi", getBodyIsNil)
	http.HandleFunc("/url/query", queryParams)
	http.HandleFunc("/wholeUrl", wholeUrl)
	http.HandleFunc("/header", header)
	http.HandleFunc("/form", form)
	http.HandleFunc("/signup", SignUp)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

/*
Postman 接口参数：
Params：传入参数值，是在url中添加问好后面的部分，例如localhost:8080/body/once?a=b，
Authorization：授权
Header：添加的HTTP的header
Body：指传过去的content Type，即请求body

Postman 回应参数：
Body：响应body数据
Headers: 相应头部信息
*/
