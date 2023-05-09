package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	jsonpatch "github.com/evanphx/json-patch"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	clientv3 "go.etcd.io/etcd/client/v3"
	"k8s.io/klog/v2"
	"practice_ctl/pkg/apimachinery/runtime"
	appsv1 "practice_ctl/pkg/apis/apps/v1"
	metav1 "practice_ctl/pkg/apis/meta"
	"practice_ctl/pkg/etcd"
	"practice_ctl/pkg/util/helpers"
)

var CarMap = map[string]runtime.Object{}

func InitCar() {
	// 初始化一个对象
	init := &appsv1.Car{
		TypeMeta: metav1.TypeMeta{
			ApiVersion: "apps/v1",
			Kind:       "Car",
		},
		ObjectMeta: metav1.ObjectMeta{
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
	strKey, strValue := parseEtcdDataCar(init)
	_ = etcd.Put(strKey, strValue)
}

func getCar(name string) (*appsv1.Car, error) {
	if car, ok := CarMap[name]; ok {
		strKey, _ := parseEtcdDataCar(car)
		err := etcd.Get(strKey)
		klog.Info("get key: ", strKey)
		if err != nil {
			klog.Errorf("get key error: ", strKey, err)
			return nil, err
		}
		return car.(*appsv1.Car), nil
	}
	return nil, nil
}

func deleteCar(name string) error {
	if c, ok := CarMap[name]; ok {
		strKey, _ := parseEtcdDataCar(c)
		car := c.(*appsv1.Car)
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

		carList.Item = append(carList.Item, v.(*appsv1.Car))
	}

	return carList, nil
}

func parseEtcdDataCar(o runtime.Object) (string, string) {
	car := o.(*appsv1.Car)
	strKey := "/" + car.Kind + "/" + car.Name
	strValue := car.Name

	return strKey, strValue
}

func createOrUpdateCar(o runtime.Object) (*appsv1.Car, error) {
	car := o.(*appsv1.Car)
	if o, ok := CarMap[car.Name]; ok {
		old := o.(*appsv1.Car)
		klog.Infof("find the apple %v, and update it!", old.Name)
		old.Name = car.Name
		old.Annotations = car.Annotations
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
		TypeMeta: metav1.TypeMeta{
			ApiVersion: "apps/v1",
			Kind:       "Car",
		},
		ObjectMeta: metav1.ObjectMeta{
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
	strKey, strValue := parseEtcdDataCar(new)
	klog.Info("create key: ", strKey)
	err := etcd.Put(strKey, strValue)
	if err != nil {
		klog.Errorf("create key error: ", strKey, err)
		return new, err
	}
	return new, nil
}

func createCar(o runtime.Object) (*appsv1.Car, error) {
	car := o.(*appsv1.Car)
	// 如果查到就抛错
	if _, ok := CarMap[car.Name]; ok {
		return nil, errors.New("this car is created ")
	}
	new := &appsv1.Car{
		TypeMeta: metav1.TypeMeta{
			ApiVersion: "apps/v1",
			Kind:       "Car",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        car.Name,
			Annotations: car.Annotations,
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

func updateCar(o runtime.Object) (*appsv1.Car, error) {
	car := o.(*appsv1.Car)
	// 重新赋值
	if o, ok := CarMap[car.Name]; ok {
		old := interface{}(o).(*appsv1.Car)
		old.Name = car.Name
		old.Annotations = car.Annotations
		old.Spec.Price = car.Spec.Price
		old.Spec.Brand = car.Spec.Brand
		old.Spec.Color = car.Spec.Color
		return old, nil
	}

	// 如果查到就抛错
	return nil, errors.New("this car is not found")
}

func patchCar(o runtime.Object) (*appsv1.Car, error) {
	car := o.(*appsv1.Car)
	// 重新赋值
	if o, ok := CarMap[car.Name]; ok {
		old := interface{}(o).(*appsv1.Car)
		car.Status = appsv1.CarStatus{
			Status: "updated",
		}
		delete(CarMap, old.Name) // delete car from map

		// 获取patch对象与原对象的差距
		patch, err := jsonpatch.CreateMergePatch(helpers.ToJson(old), helpers.ToJson(car))
		if err != nil {
			return nil, errors.New("car patch error")
		}
		newCar := &appsv1.Car{
			TypeMeta: metav1.TypeMeta{
				ApiVersion: "apps/v1",
				Kind:       "Car",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: car.Name,
			},
			Status: appsv1.CarStatus{
				Status: "updated",
			},
		}
		// 创建一个新对象给patch的差距
		nc, err := jsonpatch.MergePatch(helpers.ToJson(newCar), patch)
		if err != nil {
			return nil, errors.New("car patch create error")
		}

		var ccc appsv1.Car
		err = json.Unmarshal(nc, &ccc)
		if err != nil {
			return nil, errors.New("car patch unmarshal error")
		}
		// 记住last-apply字段
		ccc.Annotations = map[string]string{
			"last-applied-configuration": string(helpers.ToJson(ccc)),
		}
		CarMap[old.Name] = &ccc

		strKey, strValue := parseEtcdDataCar(&ccc)
		klog.Info("update key: ", strKey)
		err = etcd.Put(strKey, strValue)
		if err != nil {
			klog.Errorf("patch key error: ", strKey, err)
			return old, err
		}

		return &ccc, nil
	}

	klog.Infof("not find this car, and create it!")

	new := &appsv1.Car{
		TypeMeta: metav1.TypeMeta{
			ApiVersion: "apps/v1",
			Kind:       "Car",
		},
		ObjectMeta: metav1.ObjectMeta{
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

	CarMap[new.Name] = new
	strKey, strValue := parseEtcdDataCar(new)
	klog.Info("create key: ", strKey)
	err := etcd.Put(strKey, strValue)
	if err != nil {
		klog.Errorf("create key error: ", strKey, err)
		return new, err
	}

	// 如果查到就抛错
	return nil, errors.New("this car is not found")
}

type CarCtl struct {
}

func NewCarCtl() *CarCtl {
	return &CarCtl{}
}

func (a *CarCtl) GetCar(c *gin.Context) {
	name := c.Query("name")

	res, err := getCar(name)
	if err != nil {
		fmt.Println("get err!")
		c.JSON(400, gin.H{"error": err})
		return
	}

	c.JSON(200, res)
}

func (a *CarCtl) CreateCar(c *gin.Context) {
	var r *appsv1.Car
	if err := c.ShouldBindJSON(&r); err != nil {
		fmt.Println("bind json err!")
		c.JSON(400, gin.H{"error": err})
		return
	}
	res, err := createOrUpdateCar(r)
	if err != nil {
		fmt.Println("create err!")
		c.JSON(400, gin.H{"error": err})
		return
	}
	c.JSON(200, res)
	return

}

func (a *CarCtl) UpdateCar(c *gin.Context) {
	var r *appsv1.Car
	if err := c.ShouldBindJSON(&r); err != nil {
		fmt.Println("bind json err!")
		c.JSON(400, gin.H{"error": err})
		return
	}
	res, err := createOrUpdateCar(r)
	if err != nil {
		fmt.Println("update err!")
		c.JSON(400, gin.H{"error": err})
		return
	}

	c.JSON(200, res)
	return

}

func (a *CarCtl) DeleteCar(c *gin.Context) {
	name := c.Query("name")

	err := deleteCar(name)
	if err != nil {
		fmt.Println("get err!")
		c.JSON(400, gin.H{"error": err})
		return
	}

	c.JSON(200, gin.H{"ok": "ok"})
	return
}

func (a *CarCtl) ListCar(c *gin.Context) {

	res, err := listCar()
	if err != nil {
		fmt.Println("list err!")
		c.JSON(400, gin.H{"error": err})
		return
	}
	c.JSON(200, res)
	return

}

// 使用ws连接实现类似watch的实时传递
func (a *CarCtl) WatchCar(c *gin.Context) {
	// 升级请求
	client, err := Upgrader.Upgrade(c.Writer, c.Request, nil) //升级
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

// ws连接，用于watch操错
type WsClientCar struct {
	conn      *websocket.Conn
	writeChan chan *WatchCar // 写队列chan
	closeChan chan struct{}  // 通知停止chan
}

func NewWsClientCar(conn *websocket.Conn, writeChan chan *WatchCar, closeChan chan struct{}) *WsClientCar {
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
				w.closeChan <- struct{}{}
				break

			}

		}
	}
}

type WatchCar struct {
	Car *appsv1.Car
	// 区分事件类型 目前就是put delete
	ObjectType string
}

// watchCar 从etcd中使用watch能力，当监听到有对象put或delete时，
// watcher.ResultChan会接收到;并在内存中查找出真实对象，放入outputC中
// 从outputC中放入 ws准备写入的writeChan中
func (w *WsClientCar) watchCar(applePrefix string) {

	outputC := make(chan *WatchCar, 1000)

	watcher := etcd.Watch(applePrefix, clientv3.WithPrefix())
	for {
		select {
		case v, ok := <-watcher.ResultChan:
			if ok {
				// TODO: 可以新增事件类型：put update delete等
				for _, event := range v.Events {
					fmt.Println("value: ", string(event.Kv.Value))
					name := string(event.Kv.Value)
					// 区分事件类型
					var objectType string
					if event.Type == clientv3.EventTypePut {
						objectType = "put"
						if o, ok := CarMap[name]; ok {
							car := o.(*appsv1.Car)
							klog.Info(car.Name, car.Kind, car.Spec)
							klog.Infof("放入output中")
							res := &WatchCar{
								Car:        car,
								ObjectType: objectType,
							}
							outputC <- res
						}
					} else if event.Type == clientv3.EventTypeDelete {
						objectType = "delete"
						klog.Info("delete: ", objectType)
						res := &WatchCar{
							Car:        nil,
							ObjectType: objectType,
						}
						outputC <- res
					}

				}
			}
		case watchApple := <-outputC:
			klog.Infof("放入writeChan中")
			w.writeChan <- watchApple
		}
	}

}
