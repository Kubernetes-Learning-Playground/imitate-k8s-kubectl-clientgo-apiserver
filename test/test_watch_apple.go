package main

import (
	"encoding/json"
	"fmt"
	"k8s.io/klog/v2"
	"practice_ctl/pkg/storeapiserver/controllers"
	"practice_ctl/pkg/util/stores"
	"practice_ctl/pkg/util/stores/rest"
	"time"
)


func main() {
	//// 配置文件
	config := &rest.Config{
		Host:    fmt.Sprintf("http://localhost:8888"),
		Timeout: time.Second,
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
