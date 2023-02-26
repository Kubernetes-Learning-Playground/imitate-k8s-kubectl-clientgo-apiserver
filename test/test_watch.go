package main

import (
	"encoding/json"
	"fmt"
	"k8s.io/klog/v2"
	appsv1 "practice_ctl/pkg/apis/apps/v1"
	v1 "practice_ctl/pkg/apis/core/v1"
	"practice_ctl/pkg/util/stores"
	"practice_ctl/pkg/util/stores/rest"
	"time"
)

//
func main() {
	//// 配置文件
	config := &rest.Config{
		Host:    fmt.Sprintf("http://localhost:8080"),
		Timeout: time.Second,
	}
	clientSet := stores.NewForConfig(config)

	// watch apple对象
	res := clientSet.CoreV1().Apple().Watch()
	for i := range res.WChan {
		r := i.([]byte)
		var resApple v1.Apple
		err := json.Unmarshal(r, &resApple)
		if err != nil {
			klog.Error(err)
			return
		}
		klog.Info(resApple)
	}

	// watch car对象
	res1 := clientSet.AppsV1().Car().Watch()
	for i := range res1.WChan {
		r := i.([]byte)
		var resCar appsv1.Car
		err := json.Unmarshal(r, &resCar)
		if err != nil {
			klog.Error(err)
			return
		}
		klog.Info(resCar)
	}
}
