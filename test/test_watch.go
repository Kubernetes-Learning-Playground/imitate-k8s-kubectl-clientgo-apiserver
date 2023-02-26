package main

import (
	"encoding/json"
	"fmt"
	"k8s.io/klog/v2"
	v1 "practice_ctl/pkg/apis/apps/v1"
	"practice_ctl/pkg/util/stores"
	"practice_ctl/pkg/util/stores/rest"
	"time"
)

//
func main() {
	//// 配置文件
	//config := &rest.Config{
	//	Host:    fmt.Sprintf("http://localhost:8080"),
	//	Timeout: time.Second,
	//}
	//clientSet := stores.NewForConfig(config)
	//
	//res := clientSet.CoreV1().Apple().Watch()
	//for i := range res.WChan {
	//	r := i.([]byte)
	//	var resApple v1.Apple
	//	err := json.Unmarshal(r, &resApple)
	//	if err != nil {
	//		klog.Error(err)
	//		return
	//	}
	//	klog.Info(resApple)
	//}

	// 配置文件
	config := &rest.Config{
		Host:    fmt.Sprintf("http://localhost:8080"),
		Timeout: time.Second,
	}
	clientSet := stores.NewForConfig(config)

	res := clientSet.AppsV1().Car().Watch()
	for i := range res.WChan {
		r := i.([]byte)
		var resCar v1.Car
		err := json.Unmarshal(r, &resCar)
		if err != nil {
			klog.Error(err)
			return
		}
		klog.Info(resCar)
	}
}
