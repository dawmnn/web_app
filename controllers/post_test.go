package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreatePostHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	url := "/api/v1/post"
	r.POST(url, CreatePostHandler)

	body := `{
"community_id":1,
"title":"test",
"content":"just a test"
}`

	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader([]byte(body)))

	w := httptest.NewRecorder() //创建一个响应记录器，用于捕获 HTTP 响应
	r.ServeHTTP(w, req)         //执行 HTTP 请求，并将响应写入记录器 w

	assert.Equal(t, 200, w.Code)

	//方法一：判断响应的内容是不是按预期返回了需要登录的错误
	//	assert.Contains(t, w.Body.String(), "需要登录")

	//方法二:将响应的内容反序列化到ResponseData 然后判断字段和预期是否一致
	res := new(ResponseData) //首先创建一个新的 ResponseData 结构体实例 res，用于存储反序列化后的数据
	if err := json.Unmarshal(w.Body.Bytes(), res); err != nil {
		t.Fatalf("json.Unmarshal w.Body failed,err:%v\n", err)
	}
	assert.Equal(t, res.Code, CodeNeedLogin)
}
