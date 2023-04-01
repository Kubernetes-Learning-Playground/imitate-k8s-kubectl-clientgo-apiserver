package controllers

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/gin-gonic/gin"
	"time"
)

// TimedHandler 模拟超时handler
func TimedHandler(request *restful.Request, response *restful.Response) {
	// 获取替换之后的context 它具备了超时控制

	ctx := request.Request.Context()

	// 定义响应struct
	type responseData struct {
		status int
		body   map[string]interface{}
	}

	// 创建一个done chan表明request要完成了
	doneChan := make(chan responseData)
	// 模拟API耗时的处理
	go func() {
		time.Sleep(time.Second * 500)
		doneChan <- responseData{
			status: 200,
			body:   gin.H{"hello": "world"},
		}
	}()

	// 监听两个chan谁先到达
	select {
	// 超时
	case <-ctx.Done():
		return
	// 请求完成
	case res := <-doneChan:
		response.WriteEntity(res)
	}


}

