package controllers

import "web_app/models"

// 专门用来放接口文档用到的model
// 因为我们的接口文档返回的数据格式是一致的，但是具体的data类型不一致

type _ResponsePostList struct {
	Code    ResCode                 //业务响应的状态
	Message string                  //提示信息
	Data    []*models.ApiPostDetail //数据
}
