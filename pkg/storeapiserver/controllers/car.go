package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/shenyisyn/goft-gin/goft"
	clientv3 "go.etcd.io/etcd/client/v3"
	"k8s.io/klog/v2"
	appsv1 "practice_ctl/pkg/apis/apps/v1"
	"practice_ctl/pkg/etcd"
)

var CarMap = map[string]*appsv1.Car{}

func init() {
	// 初始化一个对象
	init := &appsv1.Car{
		ApiVersion: "apps/v1",
		Kind: 		"Car",
		Metadata: appsv1.Metadata{
			Name: "initCar",
		},
		Spec: appsv1.CarSpec{
			Brand: "initBrand",
			Price: "initPrice",
			Color: "initColor",
		},
		Status: appsv1.CarStatus{
			Status: "created",
		},
	}
	CarMap[init.Name] = init

}


func getCar(name string) (*appsv1.Car, error) {
	if car, ok := CarMap[name]; ok {
		return car, nil
	}
	return nil, nil
}

func deleteCar(name string) error {
	if car, ok := CarMap[name]; ok {
		strKey, _ := parseEtcdDataCar(car)
		klog.Info("delete key: ", strKey)
		delete(CarMap, car.Name)
		err := etcd.Delete(strKey)
		if err != nil {
			klog.Errorf("delete key error: ", strKey, err)
			return err
		}
		return nil
	}
	return nil

}

func listCar() (*appsv1.CarList, error) {
	carList := &appsv1.CarList{}
	for _, v := range CarMap {
		carList.Item = append(carList.Item, v)
	}

	return carList, nil
}

func parseEtcdDataCar(car *appsv1.Car) (string, string) {
	strKey := "/" + car.Kind + "/" + car.Name
	strValue := car.Name

	return strKey, strValue
}

func createOrUpdateCar(car *appsv1.Car) (*appsv1.Car, error) {
	if old, ok := CarMap[car.Name]; ok {
		klog.Infof("find the apple %v, and update it!", old.Name)
		old.Name = car.Name
		old.Spec.Price = car.Spec.Price
		old.Spec.Brand = car.Spec.Brand
		old.Spec.Color = car.Spec.Color
		old.Status.Status = "updated"

		strKey, strValue := parseEtcdDataCar(car)
		klog.Info("update key: ", strKey)
		err := etcd.Put(strKey, strValue)
		if err != nil {
			klog.Errorf("update key error: ", strKey, err)
			return old, err
		}
		return old, nil
	}
	klog.Infof("not find this car, and create it!")

	new := &appsv1.Car{
		ApiVersion: "apps/v1",
		Kind: "Car",
		Metadata: appsv1.Metadata{
			Name: car.Name,
		},
		Spec: appsv1.CarSpec{
			Brand: car.Spec.Brand,
			Price: car.Spec.Price,
			Color: car.Spec.Color,
		},
		Status: appsv1.CarStatus{
			Status: "created",
		},
	}


	// 存入map
	CarMap[car.Name] = new

	return new, nil
}


func createCar(car *appsv1.Car) (*appsv1.Car, error) {
	// 如果查到就抛错
	if _, ok := CarMap[car.Name]; ok {
		return nil, errors.New("this car is created ")
	}
	new := &appsv1.Car{
		ApiVersion: "apps/v1",
		Kind: "Car",
		Metadata: appsv1.Metadata{
			Name: car.Name,
		},
		Spec: appsv1.CarSpec{
			Brand: car.Spec.Brand,
			Price: car.Spec.Price,
			Color: car.Spec.Color,
		},
		Status: appsv1.CarStatus{
			Status: "created",
		},
	}


	// 存入map
	CarMap[car.Name] = new

	return new, nil

}

func updateCar(car *appsv1.Car) (*appsv1.Car, error) {
	// 重新赋值
	if old, ok := CarMap[car.Name]; ok {
		old.Name = car.Name
		old.Spec.Price = car.Spec.Price
		old.Spec.Brand = car.Spec.Brand
		old.Spec.Color = car.Spec.Color
		return old, nil
	}


	// 如果查到就抛错
	return nil, errors.New("this car is not found")


}

type CarCtl struct {
}


func (a *CarCtl) GetCar(c *gin.Context) goft.Json {
	name := c.Query("name")

	res, err := getCar(name)
	if err != nil {
		fmt.Println("get err!")
		return gin.H{"code": "400", "error": err}
	}


	return res
}

func (a *CarCtl) CreateCar(c *gin.Context) goft.Json {
	var r *appsv1.Car
	if err := c.ShouldBindJSON(&r); err != nil {
		fmt.Println("bind json err!")
		return gin.H{"code": "400", "error": err}
	}
	res, err := createOrUpdateCar(r)
	if err != nil {
		fmt.Println("create err!")
		return gin.H{"code": "400", "error": err}
	}

	return res

}

func (a *CarCtl) UpdateCar(c *gin.Context) goft.Json {
	var r *appsv1.Car
	if err := c.ShouldBindJSON(&r); err != nil {
		fmt.Println("bind json err!")
		return gin.H{"code": "400", "error": err}
	}
	res, err := createOrUpdateCar(r)
	if err != nil {
		fmt.Println("update err!")
		return gin.H{"code": "400", "error": err}
	}

	return res

}

func (a *CarCtl) DeleteCar(c *gin.Context) goft.Json {
	name := c.Query("name")

	err := deleteCar(name)
	if err != nil {
		fmt.Println("get err!")
		return gin.H{"code": "400", "error": err}
	}


	return nil
}

func (a *CarCtl) ListCar(c *gin.Context) goft.Json {

	res, err := listCar()
	if err != nil {
		fmt.Println("list err!")
		return gin.H{"code": "400", "error": err}
	}
	return res

}

// 使用ws连接实现类似watch的实时传递
func(a *CarCtl) WatchCar(c *gin.Context) (v goft.Void) {
	// 升级请求
	client, err := Upgrader.Upgrade(c.Writer,c.Request,nil)  //升级
	if err != nil {
		klog.Errorf("ws connect error", err)
		return
	}
	writeC := make(chan *appsv1.Car)
	stopC := make(chan struct{})
	ws := NewWsClientCar(client, writeC, stopC)
	// 启动两个goroutine实现
	go ws.WriteLoop()
	go ws.watchCar("/CAR")

	return
}

// ws连接，用于watch操错
type WsClientCar struct {
	conn *websocket.Conn
	writeChan chan *appsv1.Car // 写队列chan
	closeChan chan struct{}  // 通知停止chan
}

func NewWsClientCar(conn *websocket.Conn, writeChan chan *appsv1.Car, closeChan chan struct{}) *WsClientCar {
	return &WsClientCar{conn: conn, writeChan: writeChan, closeChan: closeChan}
}

// WriteLoop 发送给对应的客户端
func (w *WsClientCar) WriteLoop() {
	for {
		select {
		case msg := <-w.writeChan:
			klog.Infof("从writeChan读中")
			r, err := json.Marshal(msg)
			if err != nil {
				klog.Error(err)
			}
			klog.Infof("立即发送")
			if err := w.conn.WriteMessage(websocket.TextMessage, r); err != nil {
				klog.Errorf("发送错误")
				w.conn.Close()
				w.closeChan<- struct{}{}
				break

			}

		}
	}
}

// watchCar 从etcd中使用watch能力，当监听到有对象put或delete时，
// watcher.ResultChan会接收到;并在内存中查找出真实对象，放入outputC中
// 从outputC中放入 ws准备写入的writeChan中
func (w *WsClientCar) watchCar(applePrefix string)  {

	outputC := make(chan *appsv1.Car, 1000)

	watcher := etcd.Watch(applePrefix, clientv3.WithPrefix())
	for {
		select {
		case v, ok := <-watcher.ResultChan:
			if ok {
				// TODO: 可以新增事件类型：put update delete等
				for _, event := range v.Events {
					fmt.Println("value: ", string(event.Kv.Value))
					name := string(event.Kv.Value)
					if car, ok := CarMap[name]; ok {
						klog.Info(car.Name, car.Kind, car.Spec)
						klog.Infof("放入output中")
						outputC <-car
					}
				}
			}
		case watchApple :=  <-outputC:
			klog.Infof("放入writeChan中")
			w.writeChan <-watchApple
		}
	}

}

func (a *CarCtl) Name() string {
	return "CarCtl"
}

// 路由
func (a *CarCtl) Build(goft *goft.Goft) {
	// GET  http://localhost:8080/v1/car
	// GET  http://localhost:8080/v1/carlist
	// POST  http://localhost:8080/v1/car
	// DELETE  http://localhost:8080/v1/car
	// PUT  http://localhost:8080/v1/car
	goft.Handle("GET", "/v1/car", a.GetCar)
	goft.Handle("GET", "/v1/carlist", a.ListCar)
	goft.Handle("POST", "/v1/car", a.CreateCar)
	goft.Handle("DELETE", "/v1/car", a.DeleteCar)
	goft.Handle("PUT", "/v1/car", a.UpdateCar)
	goft.Handle("GET", "/v1/car/watch", a.WatchCar)

}

