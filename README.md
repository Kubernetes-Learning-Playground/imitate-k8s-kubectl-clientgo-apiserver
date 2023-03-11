## 模拟k8s中的kubectl命令模式与clientSet调用的练习。

### 目前还在一边开发一边设计，主要为了学习k8s的设计理念而创建的。

### api server资源
![](https://github.com/googs1025/imitate-k8s-kubectl-clientSet/blob/main/img/%E6%B5%81%E7%A8%8B%E5%9B%BE1.jpg?raw=true)

模仿k8s api server 的资源对象分类，分为apps、core等。
1. 目前支持两种资源大类，并可以在api server中不断扩展
    a. core：可以看成水果类资源
    b. apps：可以看成汽车类资源
2. 提供每种资源的create update delete get list 方法，目前支持list create apply等命令行操作
3. 资源对象改成声明式api的形式，每次更新底层都是使用createOrUpdate方法
### RoadMap 
**TODO** 未来提供scheme注册表。
**TODO** 新增aggregator apiserver功能
**TODO** 实现informer机制。(时间周期长，预计缓更)
### clientSet 风格客户端封装
如下图所示：基于net/http基础库的封装，并依据k8s风格封装http CRUD接口。

![](https://github.com/googs1025/imitate-k8s-kubectl-clientSet/blob/main/img/%E6%B5%81%E7%A8%8B%E5%9B%BE.jpg?raw=true)

功能：目前实现apple car资源对象(ex: pod)，并实现**GET LIST DELETE CREATE UPDATE WATCH** 方法

#### 范例文件
```bigquery
    // 创建操作
	a := &v1.Apple{
		ApiVersion: "core/v1",
		Kind: "APPLE",
		Metadata: v1.Metadata{
			Name: "apple1",
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
	c, err := clientSet.CoreV1().Apple().Create(a)
    if err != nil {
		fmt.Println(err)
	}
	fmt.Println("name:", c.Name,  "size:", c.Spec.Size, "color:", c.Spec.Color, "place:", c.Spec.Place, "price:", c.Spec.Price)

	apple1, err := clientSet.CoreV1().Apple().Get(c.Name)
    if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("name: ", apple1.Name)

	aaa := &v1.Apple{
		ApiVersion: "core/v1",
		Kind: "APPLE",
		Metadata: v1.Metadata{
			Name: "apple1",
		},
		Spec: v1.AppleSpec{
			Size: "apple1dddd",
			Color: "apple1ccc",
			Place: "apple1ccc",
			Price: "apple1ccc",
		},
	}

	appleUpdate, err := clientSet.CoreV1().Apple().Update(aaa)
    if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("name: ", appleUpdate.Name,  "size: ", appleUpdate.Spec.Size, "color: ", appleUpdate.Spec.Color, "place: ", appleUpdate.Spec.Place, "price: ", appleUpdate.Spec.Price)

	appleList, err := clientSet.CoreV1().Apple().List()
    if err != nil {
		log.Fatalln(err)
	}
	for _, apple := range appleList.Item {
		fmt.Println(apple.Name)
	}
```

#### watch 操作
```
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
```

### kubectl 风格命令行封装
![](https://github.com/googs1025/imitate-k8s-kubectl-clientSet/blob/main/img/%E6%B5%81%E7%A8%8B%E5%9B%BE11.jpg?raw=true)
底层使用client-go对资源进行CRUD操作。
#### 使用方法
1. 还没编译，所以执行使用 go run storectl.go xxxx 测试，编译后放入执行的/bin目录中，可以使用storectl list xxxx
2. 目前仅支持list create apply delete等操作。
3. 同时支持json yaml两种方式创建与修改。
```bigquery
# json yaml 范例文件 json目录 yaml目录
# 进入：cmd/storectl目录中
# 获取命令： storectl list apples 
➜  storectl git:(main) ✗ go run storectl.go apply cars ../../json/car.json
+-----------+-----------+-----------+-----------+----------+
| APPLE名称 |   PRICE   |   PLACE   |   COLOR   |   SIZE   |
+-----------+-----------+-----------+-----------+----------+
| aaaa      | ccc       | ccc       | aaa       | aaa      |
| initApple | initPrice | initPlace | initColor | initSize |
+-----------+-----------+-----------+-----------+----------+
# 创建命令 ： 类似 storectl create apples aaa.json
➜  storectl git:(main) ✗ go run storectl.go create apples aaa.json
name:dddafjjhadsklfhaaaa is created

```

### 项目目录

```bigquery
├── README.md
├── cmd # 编译main文件
├── pkg
│   ├── apis # 用来存放资源对象
│   │   └── core
│   │       ├── unverstioned
│   │       │   └── version.go
│   │       └── v1
│   │           ├── apple.go
│   │           └── version.go
│   ├── storeapiserver
│   │   ├── auth
│   │   ├── configs
│   │   └── controllers # api server控制器
│   │       ├── apple.go
│   │       └── version.go
│   ├── storectl
│   │   ├── cmd 
│   │   │   ├── base.go
│   │   │   └── versionCmd.go
│   │   └── config # 客户端配置文件
│   │       └── config.go
│   └── util
│       ├── helpers
│       │   └── filehelper.go
│       └── stores
│           ├── clientset.go
│           ├── rest # RESTClient http库封装
│           │   ├── client.go
│           │   ├── config.go
│           │   ├── request.go
│           │   └── result.go
│           └── typed 
│               └── core # core资源对象的 client 封装
│                   ├── apple.go
│                   ├── core_client.go
│                   └── version.go
└── test.go
```
