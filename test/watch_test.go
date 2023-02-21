package test

import (
	"fmt"
	v1 "practice_ctl/pkg/apis/core/v1"
	"practice_ctl/pkg/storeapiserver/controllers"
	"practice_ctl/pkg/util/stores"
	"practice_ctl/pkg/util/stores/rest"
	"testing"
	"time"
)

func TestWatch(t *testing.T) {
	// 配置文件
	config := &rest.Config{
		Host:    fmt.Sprintf("http://localhost:8080"),
		Timeout: time.Second,
	}
	clientSet := stores.NewForConfig(config)

	// 创建操作
	a := &v1.Apple{
		ApiVersion: "core/v1",
		Kind: "APPLE",
		Metadata: v1.Metadata{
			Name: "apple2111",
		},
		Spec: v1.AppleSpec{
			Size: "apple1",
			Color: "apple1",
			Place: "apple1",
			Price: "apple1",
		},
		Status: v1.AppleStatus{
			Status: "created",
		},

	}
	_, err := clientSet.CoreV1().Apple().Create(a)
	if err != nil {
		fmt.Println(err)
	}
	cc := controllers.AppleCtl{}

	cc.WatchApple()

}
