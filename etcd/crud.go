package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
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

// Watch 监听
func Watch(key string, opts ...clientv3.OpOption) {
	ch := Cli.Watch(ctx, key, opts...)
	for v := range ch {
		for _, val := range v.Events {
			PrintJSON(val)
		}
	}
}

func PrintJSON(v interface{}) {
	b, _ := json.Marshal(v)
	fmt.Printf("%s\n", b)
}
