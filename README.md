## 模拟k8s中的kubectl命令模式与clientSet调用的练习。

### clientSet 风格客户端封装
如下图所示：基于net/http基础库的封装，并依据k8s风格封装http CRUD接口。

![](https://github.com/googs1025/imitate-k8s-kubectl-clientSet/blob/main/img/%E6%B5%81%E7%A8%8B%E5%9B%BE.jpg?raw=true)

功能：目前实现apple资源对象(ex: pod)，并实现**GET LIST DELETE CREATE UPDATE** 方法

#### 范例文件
```bigquery
    // 创建操作
	a := &v1.Apple{
		Name: "apple1",
		Size: "apple1",
		Color: "apple1",
		Place: "apple1",
		Price: "apple1",
	}
	c, err := clientSet.Core().Apple().Create(a)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("name:", c.Name,  "size:", c.Size, "color:", c.Color, "place:", c.Place, "price:", c.Price)

	apple1, err := clientSet.Core().Apple().Get(c.Name)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("name: ", apple1.Name)

	aaa := &v1.Apple{
		Name: "apple1",
		Size: "apple1dddd",
		Color: "apple1ccc",
		Place: "apple1ccc",
		Price: "apple1ccc",
	}

	appleupdate, err := clientSet.Core().Apple().Update(aaa)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("name: ", appleupdate.Name,  "size: ", appleupdate.Size, "color: ", appleupdate.Color, "place: ", appleupdate.Place, "price: ", appleupdate.Price)

	appleList, err := clientSet.Core().Apple().List()
	if err != nil {
		log.Fatalln(err)
	}
	for _, apple := range appleList.Item {
		fmt.Println(apple.Name)
	}
```


### kubectl 风格命令行封装