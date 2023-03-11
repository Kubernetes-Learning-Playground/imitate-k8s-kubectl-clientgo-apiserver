package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	clientv3 "go.etcd.io/etcd/client/v3"
	"k8s.io/klog/v2"
	v1 "practice_ctl/pkg/apis/core/v1"
	"practice_ctl/pkg/etcd"
)

var AppleMap = map[string]*v1.Apple{}

func init() {
	// 初始化一个对象
	init := &v1.Apple{
		ApiVersion: "core/v1",
		Kind: "Apple",
		Metadata: v1.Metadata{
			Name: "initApple",
		},
		Spec: v1.AppleSpec{
			Place: "initPlace",
			Price: "initPrice",
			Size:  "initSize",
			Color: "initColor",
		},
		Status: v1.AppleStatus{
			Status: "created",
		},

	}

	AppleMap[init.Name] = init
	strKey, strValue := parseEtcdData(init)
	_ = etcd.Put(strKey, strValue)


}


func getApple(name string) (*v1.Apple, error) {
	if apple, ok := AppleMap[name]; ok {

		strKey, _ := parseEtcdData(apple)
		err := etcd.Get(strKey)
		klog.Info("get key: ", strKey)
		if err != nil {
			klog.Errorf("get key error: ", strKey, err)
			return apple, nil
		}

		return apple, nil
	}
	return nil, nil
}

func deleteApple(name string) error {
	if apple, ok := AppleMap[name]; ok {
		strKey, _ := parseEtcdData(apple)
		klog.Info("delete key: ", strKey)
		err := etcd.Delete(strKey)
		delete(AppleMap, apple.Name)
		if err != nil {
			klog.Errorf("delete key error: ", strKey, err)
			return err
		}
		return nil
	}
	return nil

}

func listApple() (*v1.AppleList, error) {
	appleList := &v1.AppleList{}
	for _, v := range AppleMap {
		appleList.Item = append(appleList.Item, v)
	}

	return appleList, nil
}

func createOrUpdateApple(apple *v1.Apple) (*v1.Apple, error) {
	if old, ok := AppleMap[apple.Name]; ok {
		klog.Infof("find the apple %v, and update it!", old.Name)
		old.Name = apple.Name
		old.Spec.Price = apple.Spec.Price
		old.Spec.Place = apple.Spec.Place
		old.Spec.Size = apple.Spec.Size
		old.Spec.Color = apple.Spec.Color
		old.Status.Status = "updated"

		strKey, strValue := parseEtcdData(apple)
		klog.Info("update key: ", strKey)
		err := etcd.Put(strKey, strValue)
		if err != nil {
			klog.Errorf("update key error: ", strKey, err)
			return old, err
		}

		return old, nil
	}
	klog.Infof("not find this apple, and create it!")
	a := v1.AppleStatus{
		Status: "created",
	}
	new := &v1.Apple{
		ApiVersion: apple.ApiVersion,
		Kind: apple.Kind,
		Metadata: v1.Metadata{
			Name: apple.Name,
		},
		Spec: apple.Spec,
		Status: a,
	}


	// 存入map
	AppleMap[apple.Name] = new

	strKey, strValue := parseEtcdData(new)
	klog.Info("create key: ", strKey)
	err := etcd.Put(strKey, strValue)
	if err != nil {
		klog.Errorf("create key error: ", strKey, err)
		return new, err
	}

	return new, nil
}

// ws连接，用于watch操错
type WsClientApple struct {
	conn *websocket.Conn
	writeChan chan *WatchApple // 写队列chan
	closeChan chan struct{}  // 通知停止chan
}

func NewWsClientApple(conn *websocket.Conn, writeChan chan *WatchApple, closeChan chan struct{}) *WsClientApple {
	return &WsClientApple{conn: conn, writeChan: writeChan, closeChan: closeChan}
}

// WriteLoop 发送给对应的客户端
func (w *WsClientApple) WriteLoop() {
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

type WatchApple struct {
	Apple  *v1.Apple
	// 区分事件类型 目前就是put delete
	ObjectType string
}


// watchApple 从etcd中使用watch能力，当监听到有对象put或delete时，
// watcher.ResultChan会接收到;并在内存中查找出真实对象，放入outputC中
// 从outputC中放入 ws准备写入的writeChan中
func (w *WsClientApple) watchApple(applePrefix string)  {

	outputC := make(chan *WatchApple, 1000)

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
						if apple, ok := AppleMap[name]; ok {
							klog.Info(apple.Name, apple.Kind, apple.Spec)
							klog.Infof("放入output中")
							res := &WatchApple{
								Apple: apple,
								ObjectType: objectType,
							}
							outputC <-res
						}
					} else if event.Type == clientv3.EventTypeDelete {
						objectType = "delete"
						res := &WatchApple{
							Apple: nil,
							ObjectType: objectType,
						}
						outputC <-res
					}


				}
			} 
		case watchApple :=  <-outputC:
			klog.Infof("放入writeChan中")
			w.writeChan <-watchApple
		}
	}

}

//func createApple(apple *v1.Apple) (*v1.Apple, error) {
//	// 如果查到就抛错
//	if _, ok := AppleMap[apple.Name]; ok {
//		return nil, errors.New("this apple is created ")
//	}
//	new := &v1.Apple{
//		Name: apple.Name,
//		Size: apple.Size,
//		Price: apple.Price,
//		Place: apple.Place,
//		Color: apple.Color,
//	}
//
//	// 存入map
//	AppleMap[apple.Name] = new
//
//	return new, nil
//
//}

//func updateApple(apple *v1.Apple) (*v1.Apple, error) {
//	// 重新赋值
//	if old, ok := AppleMap[apple.Name]; ok {
//		old.Name = apple.Name
//		old.Price = apple.Price
//		old.Place = apple.Place
//		old.Size = apple.Size
//		old.Color = apple.Color
//		return old, nil
//	}
//
//
//	// 如果查到就抛错
//	return nil, errors.New("this apple is not found ")
//
//
//}

type AppleCtl struct {
}

func NewAppleCtl() *AppleCtl {
	return &AppleCtl{}
}


func (a *AppleCtl) GetApple(c *gin.Context) {
	name := c.Query("name")

	res, err := getApple(name)
	if err != nil {
		fmt.Println("get err!")
		c.JSON(400, gin.H{"error": err})
		return
	}

	c.JSON(200, res)
}

func (a *AppleCtl) CreateApple(c *gin.Context) {
	var r *v1.Apple
	if err := c.ShouldBindJSON(&r); err != nil {
		fmt.Println("bind json err!")
		c.JSON(400, gin.H{"error": err})
		return
	}
	res, err := createOrUpdateApple(r)
	if err != nil {
		fmt.Println("create err!")
		c.JSON(400, gin.H{"error": err})
		return
	}

	c.JSON(200, res)
	return
}

func (a *AppleCtl) UpdateApple(c *gin.Context) {
	var r *v1.Apple
	if err := c.ShouldBindJSON(&r); err != nil {
		fmt.Println("bind json err!")
		c.JSON(400, gin.H{"error": err})
		return
	}
	res, err := createOrUpdateApple(r)
	if err != nil {
		fmt.Println("update err!")
		c.JSON(400, gin.H{"error": err})
		return
	}

	c.JSON(200, res)
	return

}

func (a *AppleCtl) DeleteApple(c *gin.Context) {
	name := c.Query("name")

	err := deleteApple(name)
	if err != nil {
		fmt.Println("get err!")
		c.JSON(400, gin.H{"error": err})
		return
	}

	c.JSON(200, gin.H{"ok": "ok"})
	return
}

func (a *AppleCtl) ListApple(c *gin.Context) {

	res, err := listApple()
	if err != nil {
		fmt.Println("list err!")
		c.JSON(400, gin.H{"error": err})
		return
	}
	c.JSON(200, res)
	return

}


func parseEtcdData(apple *v1.Apple) (string, string) {
	strKey := "/" + apple.Kind + "/" + apple.Name
	strValue := apple.Name

	return strKey, strValue
}

// 使用ws连接实现类似watch的实时传递
func(a *AppleCtl) WatchApple(c *gin.Context) {
	// 升级请求
	client, err := Upgrader.Upgrade(c.Writer,c.Request,nil)  //升级
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

