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
		Token:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2ODM3MjY5NjcsInVzZXJuYW1lIjoidGVzdCJ9.wVXVKNGSYpA_c--um8ig9xn5-1svscvTUkI_js0WWFE",
	}
	clientSet := stores.NewForConfig(config)

	// watch apple对象
	res := clientSet.CoreV1().Apple().Watch()
	for i := range res.WChan {
		r := i.([]byte)
		klog.Info(string(r))
		var resApple controllers.WatchApple
		err := json.Unmarshal(r, &resApple)
		if err != nil {
			klog.Error(err)
			return
		}
		klog.Info("res: ", resApple.Apple, "event type: ", resApple.ObjectType)
	}

}
