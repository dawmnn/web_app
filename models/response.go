package models

// ResponseSuccess 代表成功响应
type ResponseSuccess struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// ResponseError 代表错误响应
type ResponseError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}
