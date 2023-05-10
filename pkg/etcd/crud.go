package etcd

import (
	"context"
	"encoding/json"
	clientv3 "go.etcd.io/etcd/client/v3"
	"k8s.io/klog/v2"
)

var ctx = context.TODO()

// PutAndResourceVersion 设置值与返回版本号
func PutAndResourceVersion(key, val string, opts ...clientv3.OpOption) (int, error) {
	res, err := Cli.Put(ctx, key, val, opts...)
	if err != nil {
		return 0, err
	}
	resourceVersion := PrintJSON(res)
	return resourceVersion, nil
}

// Get 获取值
func Get(key string, opts ...clientv3.OpOption) error {
	res, err := Cli.Get(ctx, key, opts...)
	if err != nil {
		return err
	}
	PrintJSON(res)
	return nil
}

// Delete 删除键值对
func Delete(key string, opts ...clientv3.OpOption) error {
	res, err := Cli.Delete(ctx, key, opts...)
	if err != nil {
		return err
	}
	PrintJSON(res)
	return nil
}

type Watcher struct {
	ResultChan WatchChan
}

func NewWatcher(resultChan WatchChan) *Watcher {
	return &Watcher{ResultChan: resultChan}
}

type WatchChan <-chan clientv3.WatchResponse

// Watch 监听
func Watch(key string, opts ...clientv3.OpOption) *Watcher {

	ch := Cli.Watch(ctx, key, opts...)
	watcher := NewWatcher(WatchChan(ch))
	return watcher
}

type ResourceVersion struct {
	Header h `json:"header"`
}

type h struct {
	//ClusterId string `json:"cluster_id"`
	Revision int `json:"revision"`
}

// PrintJSON 打印log后返回etcd中的版本号
func PrintJSON(v interface{}) int {

	var rr ResourceVersion
	b, _ := json.Marshal(v)
	err := json.Unmarshal(b, &rr)
	if err != nil {
		klog.Errorf("unmarshal error: %s", err)
		return 0
	}

	klog.Infof("do something in etcd: %s\n", b)
	return rr.Header.Revision
}
