package controllers

import (
	"fmt"
	"github.com/emicklei/go-restful/v3"
	"k8s.io/klog/v2"
	appsv1 "practice_ctl/pkg/apis/apps/v1"
)

type CarRestfulCtl struct {
}

func NewCarRestfulCtl() *CarRestfulCtl {
	return &CarRestfulCtl{}
}

func (c *CarRestfulCtl) GetCar(request *restful.Request, response *restful.Response) {
	name := request.QueryParameter("name")

	res, err := getCar(name)
	if err != nil {
		fmt.Println("get err!")
		errResp := struct {
			Code int    `json:"code"`
			Err  string `json:"err"`
		}{Code: 400, Err: err.Error()}
		response.WriteEntity(&errResp)
		return
	}

	response.WriteEntity(&res)

}


func (c *CarRestfulCtl) CreateCar(request *restful.Request, response *restful.Response) {
	var r *appsv1.Car
	if err := request.ReadEntity(&r); err != nil {
		fmt.Println("bind json err!")
		errResp := struct {
			Code int    `json:"code"`
			Err  string `json:"err"`
		}{Code: 400, Err: err.Error()}
		response.WriteEntity(&errResp)
		return
	}
	res, err := createOrUpdateCar(r)
	if err != nil {
		fmt.Println("create err!")
		errResp := struct {
			Code int    `json:"code"`
			Err  string `json:"err"`
		}{Code: 400, Err: err.Error()}
		response.WriteEntity(&errResp)
		return
	}

	response.WriteEntity(&res)

	return
}

func (c *CarRestfulCtl) UpdateCar(request *restful.Request, response *restful.Response) {
	var r *appsv1.Car
	if err := request.ReadEntity(&r); err != nil {
		fmt.Println("bind json err!")
		errResp := struct {
			Code int    `json:"code"`
			Err  string `json:"err"`
		}{Code: 400, Err: err.Error()}
		response.WriteEntity(&errResp)
		return
	}
	res, err := createOrUpdateCar(r)
	if err != nil {
		fmt.Println("create err!")
		errResp := struct {
			Code int    `json:"code"`
			Err  string `json:"err"`
		}{Code: 400, Err: err.Error()}
		response.WriteEntity(&errResp)
		return
	}


	response.WriteEntity(&res)
	return

}

func (c *CarRestfulCtl) DeleteCar(request *restful.Request, response *restful.Response) {
	name := request.QueryParameter("name")

	err := deleteCar(name)
	if err != nil {
		fmt.Println("get err!")
		errResp := struct {
			Code int    `json:"code"`
			Err  string `json:"err"`
		}{Code: 400, Err: err.Error()}
		response.WriteEntity(&errResp)
		return
	}
	resp := struct {
		Code int    	 `json:"code"`
		Ok   interface{} `json:"ok"`
	}{Code: 200, Ok: "ok"}

	response.WriteEntity(&resp)

	return
}

func (c *CarRestfulCtl) ListCar(request *restful.Request, response *restful.Response) {

	res, err := listCar()
	if err != nil {
		fmt.Println("list err!")
		errResp := struct {
			Code int    `json:"code"`
			Err  string `json:"err"`
		}{Code: 400, Err: err.Error()}
		response.WriteEntity(&errResp)
		return
	}

	resp := struct {
		Code int    	 `json:"code"`
		Item interface{} `json:"Item"`
	}{Code: 200, Item: res.Item}

	response.WriteEntity(&resp)
	return

}




// 使用ws连接实现类似watch的实时传递
func(c *CarRestfulCtl) WatchCar(request *restful.Request, response *restful.Response) {
	// 升级请求

	client, err := Upgrader.Upgrade(response, request.Request,nil)  //升级
	if err != nil {
		klog.Errorf("ws connect error", err)
		return
	}
	writeC := make(chan *WatchCar)
	stopC := make(chan struct{})
	ws := NewWsClientCar(client, writeC, stopC)
	// 启动两个goroutine实现
	go ws.WriteLoop()
	go ws.watchCar("/Car")

	return
}

