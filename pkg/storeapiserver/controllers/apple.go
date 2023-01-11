package controllers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shenyisyn/goft-gin/goft"
	"k8s.io/klog/v2"
	v1 "practice_ctl/pkg/apis/core/v1"
)

var AppleMap = map[string]*v1.Apple{}

func init() {
	// 初始化一个对象
	init := &v1.Apple{
		ApiVersion: "core/v1",
		Kind: "APPLE",
		Metadata: v1.Metadata{
			Name: "initApple",
		},
		Spec: v1.AppleSpec{
			Place: "initPlace",
			Price: "initPrice",
			Size: "initSize",
			Color: "initColor",
		},
		Status: v1.AppleStatus{
			Status: "created",
		},

	}
	AppleMap[init.Name] = init

}


func getApple(name string) (*v1.Apple, error) {
	if apple, ok := AppleMap[name]; ok {
		return apple, nil
	}
	return nil, errors.New("not found apple")
}

func deleteApple(name string) error {
	if apple, ok := AppleMap[name]; ok {
		delete(AppleMap, apple.Name)
		return nil
	}
	return errors.New("not found this apple")

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

	return new, nil
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


func (a *AppleCtl) GetApple(c *gin.Context) goft.Json {
	name := c.Query("name")

	res, err := getApple(name)
	if err != nil {
		fmt.Println("get err!")
		return gin.H{"code": "400", "error": err}
	}


	return res
}

func (a *AppleCtl) CreateApple(c *gin.Context) goft.Json {
	var r *v1.Apple
	if err := c.ShouldBindJSON(&r); err != nil {
		fmt.Println("bind json err!")
		return gin.H{"code": "400", "error": err}
	}
	res, err := createOrUpdateApple(r)
	if err != nil {
		fmt.Println("create err!")
		return gin.H{"code": "400", "error": err}
	}

	return res

}

func (a *AppleCtl) UpdateApple(c *gin.Context) goft.Json {
	var r *v1.Apple
	if err := c.ShouldBindJSON(&r); err != nil {
		fmt.Println("bind json err!")
		return gin.H{"code": "400", "error": err}
	}
	res, err := createOrUpdateApple(r)
	if err != nil {
		fmt.Println("update err!")
		return gin.H{"code": "400", "error": err}
	}

	return res

}

func (a *AppleCtl) DeleteApple(c *gin.Context) goft.Json {
	name := c.Query("name")

	err := deleteApple(name)
	if err != nil {
		fmt.Println("get err!")
		return gin.H{"code": "400", "error": err}
	}


	return nil
}

func (a *AppleCtl) ListApple(c *gin.Context) goft.Json {

	res, err := listApple()
	if err != nil {
		fmt.Println("list err!")
		return gin.H{"code": "400", "error": err}
	}
	return res

}

func (a *AppleCtl) Name() string {
	return "AppleCtl"
}

// 路由
func (a *AppleCtl) Build(goft *goft.Goft) {
	// GET  http://localhost:8080/v1/apple
	// GET  http://localhost:8080/v1/applelist
	// POST  http://localhost:8080/v1/apple
	// DELETE  http://localhost:8080/v1/apple
	// PUT  http://localhost:8080/v1/apple
	goft.Handle("GET", "/v1/apple", a.GetApple)
	goft.Handle("GET", "/v1/applelist", a.ListApple)
	goft.Handle("POST", "/v1/apple", a.CreateApple)
	goft.Handle("DELETE", "/v1/apple", a.DeleteApple)
	goft.Handle("PUT", "/v1/apple", a.UpdateApple)

}
