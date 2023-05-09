package etcd

import (
	clientv3 "go.etcd.io/etcd/client/v3"
	"k8s.io/klog/v2"
)

// Cli Etcd客户端
var Cli *clientv3.Client

// InitOrDie 初始化如果不成功，直接panic
func InitOrDie(etcdConfig clientv3.Config) {
	Cli = GetEtcdClient(etcdConfig)
}

// GetEtcdClient 获取Etcd客户端
func GetEtcdClient(cfg clientv3.Config) *clientv3.Client {
	cli, err := clientv3.New(cfg)
	if err != nil {
		klog.Error("etcd error: ", err)
		panic(err)
	}
	return cli
}

/* etcd image启动命令
rm -rf /tmp/etcd-data.tmp && mkdir -p /tmp/etcd-data.tmp && \
  docker rmi gcr.io/etcd-development/etcd:v3.4.23 || true && \
  docker run \
  -p 2379:2379 \
  -p 2380:2380 \
  --mount type=bind,source=/tmp/etcd-data.tmp,destination=/etcd-data \
  --name etcd-gcr-v3.4.23 \
  gcr.io/etcd-development/etcd:v3.4.22 \
  /usr/local/bin/etcd \
  --name s1 \
  --data-dir /etcd-data \
  --listen-client-urls http://0.0.0.0:2379 \
  --advertise-client-urls http://0.0.0.0:2379 \
  --listen-peer-urls http://0.0.0.0:2380 \
  --initial-advertise-peer-urls http://0.0.0.0:2380 \
  --initial-cluster s1=http://0.0.0.0:2380 \
  --initial-cluster-token tkn \
  --initial-cluster-state new \
  --log-level info \
  --logger zap \
  --log-outputs stderr
*/
