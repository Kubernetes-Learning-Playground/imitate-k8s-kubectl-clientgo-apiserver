package main

import (
	"encoding/json"
	"fmt"
	"k8s.io/klog/v2"
	"practice_ctl/pkg/apiserver/controllers"
	"practice_ctl/pkg/util/stores"
	"practice_ctl/pkg/util/stores/rest"
	"time"
)

func main() {
	//// 配置文件
	config := &rest.Config{
		Host:    fmt.Sprintf("http://localhost:8080"),
		Timeout: time.Second,
		Token:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2ODM3MjM4NTksInVzZXJuYW1lIjoiYWRtaW4ifQ.nGrktCEZsdbyqWSsSK80af6ncihavBB41WLm5eqpmDc",
	}
	clientSet := stores.NewForConfig(config)

	// watch car对象
	res1 := clientSet.AppsV1().Car().Watch()
	for i := range res1.WChan {
		r := i.([]byte)
		klog.Info(string(r))
		var resCar controllers.WatchCar
		err := json.Unmarshal(r, &resCar)
		if err != nil {
			klog.Error(err)
			return
		}
		klog.Info(resCar.Car, resCar.ObjectType)

	}
}
