package controllers

import (
	"fmt"
	"github.com/emicklei/go-restful/v3"
	"k8s.io/klog/v2"
	"net/http"
	v1 "practice_ctl/pkg/apis/core/v1"
)

type AppleRestfulCtl struct {
}

func NewAppleRestfulCtl() *AppleRestfulCtl {
	return &AppleRestfulCtl{}
}

func (a *AppleRestfulCtl) GetApple(request *restful.Request, response *restful.Response) {
	name := request.QueryParameter("name")

	res, err := getApple(name)
	if err != nil {
		fmt.Println("get err!")
		errResp := struct {
			Code int    `json:"code"`
			Err  string `json:"err"`
		}{Code: http.StatusBadRequest, Err: err.Error()}
		response.WriteEntity(&errResp)
		return
	}



	response.WriteEntity(&res)

}


func (a *AppleRestfulCtl) CreateApple(request *restful.Request, response *restful.Response) {
	var r *v1.Apple
	if err := request.ReadEntity(&r); err != nil {
		fmt.Println("bind json err!")
		errResp := struct {
			Code int    `json:"code"`
			Err  string `json:"err"`
		}{Code: http.StatusBadRequest, Err: err.Error()}
		response.WriteEntity(&errResp)
		return
	}
	res, err := createOrUpdateApple(r)
	if err != nil {
		fmt.Println("create err!")
		errResp := struct {
			Code int    `json:"code"`
			Err  string `json:"err"`
		}{Code: http.StatusBadRequest, Err: err.Error()}
		response.WriteEntity(&errResp)
		return
	}


	response.WriteEntity(&res)

	return
}

func (a *AppleRestfulCtl) UpdateApple(request *restful.Request, response *restful.Response) {
	var r *v1.Apple
	if err := request.ReadEntity(&r); err != nil {
		fmt.Println("bind json err!")
		errResp := struct {
			Code int    `json:"code"`
			Err  string `json:"err"`
		}{Code: http.StatusBadRequest, Err: err.Error()}
		response.WriteEntity(&errResp)
		return
	}
	res, err := createOrUpdateApple(r)
	if err != nil {
		fmt.Println("create err!")
		errResp := struct {
			Code int    `json:"code"`
			Err  string `json:"err"`
		}{Code: http.StatusBadRequest, Err: err.Error()}
		response.WriteEntity(&errResp)
		return
	}


	response.WriteEntity(&res)
	return

}

func (a *AppleRestfulCtl) DeleteApple(request *restful.Request, response *restful.Response) {
	name := request.QueryParameter("name")

	err := deleteApple(name)
	if err != nil {
		fmt.Println("get err!")
		errResp := struct {
			Code int    `json:"code"`
			Err  string `json:"err"`
		}{Code: http.StatusBadRequest, Err: err.Error()}
		response.WriteEntity(&errResp)
		return
	}
	resp := struct {
		Code int    	 `json:"code"`
		Ok   interface{} `json:"ok"`
	}{Code: http.StatusOK, Ok: "ok"}

	response.WriteEntity(&resp)

	return
}

func (a *AppleRestfulCtl) ListApple(request *restful.Request, response *restful.Response) {

	res, err := listApple()
	if err != nil {
		fmt.Println("list err!")
		errResp := struct {
			Code int    `json:"code"`
			Err  string `json:"err"`
		}{Code: http.StatusBadRequest, Err: err.Error()}
		response.WriteEntity(&errResp)
		return
	}

	resp := struct {
		Code int    	 `json:"code"`
		Item interface{} `json:"Item"`
	}{Code: http.StatusOK, Item: res.Item}

	response.WriteEntity(&resp)
	return

}


func (a *AppleRestfulCtl) PatchApple(request *restful.Request, response *restful.Response) {
	var r *v1.Apple
	if err := request.ReadEntity(&r); err != nil {
		fmt.Println("bind json err!")
		errResp := struct {
			Code int    `json:"code"`
			Err  string `json:"err"`
		}{Code: http.StatusBadRequest, Err: err.Error()}
		response.WriteEntity(&errResp)
		return
	}
	res, err := patchApple(r)
	if err != nil {
		fmt.Println("create err!")
		errResp := struct {
			Code int    `json:"code"`
			Err  string `json:"err"`
		}{Code: http.StatusBadRequest, Err: err.Error()}
		response.WriteEntity(&errResp)
		return
	}

	response.WriteEntity(&res)
	return

}


// 使用ws连接实现类似watch的实时传递
func(a *AppleRestfulCtl) WatchApple(request *restful.Request, response *restful.Response) {
	// 升级请求

	client, err := Upgrader.Upgrade(response, request.Request,nil)  //升级
	if err != nil {
		klog.Errorf("ws connect error", err)
		return
	}
	writeC := make(chan *WatchApple)
	stopC := make(chan struct{})
	ws := NewWsClientApple(client, writeC, stopC)
	// 启动两个goroutine实现
	go ws.WriteLoop()
	go ws.watchApple("/Apple")

	return
}

