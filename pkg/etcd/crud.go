package etcd

import (
	"context"
	"encoding/json"
	clientv3 "go.etcd.io/etcd/client/v3"
	"k8s.io/klog/v2"
)

var ctx = context.TODO()

// Put 设置值
func Put(key, val string, opts ...clientv3.OpOption) error {
	res, err := Cli.Put(ctx, key, val, opts...)
	if err != nil {
		return err
	}
	PrintJSON(res)
	return nil
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

func PrintJSON(v interface{}) {
	b, _ := json.Marshal(v)
	klog.Infof("do something in etcd: %s\n", b)
}
