package rest

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/websocket"
	"io"
	"k8s.io/klog/v2"
	"net/http"
	"net/url"
	appsv1 "practice_ctl/pkg/apis/apps/v1"
	v1 "practice_ctl/pkg/apis/core/v1"
	"strings"
)

// 模仿 k8s
type Request struct {
	c       *RESTClient
	req 	*http.Request

	url     url.URL
	ws      *websocket.Conn
	// TODO: 这里需要修改
	WChan   chan interface{}
}

type WsChan struct {
	Object interface{}
	Type   string
}

func NewRequest(c *RESTClient) *Request {
	req, _ := http.NewRequest("", "", nil)
	return &Request{
		c: c,
		req: req,
		WChan: make(chan interface{}, 10),

	} //默认是根
}


// Do 需要对底层库进行封装，不要暴露
// 最终执行 http请求
func (r *Request) Do() Result {
	var ret Result

	rsp, err := r.c.Do(r.req)
	if err != nil {
		ret.err = err
	} else {
		ret.rsp = rsp
	}

	return ret
}

//  method 设置
func (r *Request) Verb(verb string) *Request {
	r.req.Method = verb
	return r
}

func (r *Request) Path(p string) *Request {
	str := strings.Split(r.c.BasePath, "://")
	r.req.URL.Scheme = str[0]
	r.req.URL.Host = str[1]

	r.req.URL.Path = p
	return r
}

func (r *Request) GetCarName(name string) *Request {

	q := r.req.URL.Query()

	queryParam := map[string]string{
		"name": name,
	}
	for k, v := range queryParam {
		q.Add(k, v)
	}
	r.req.URL.RawQuery = q.Encode()

	return r
}

func (r *Request) CreateCar(car *appsv1.Car) *Request {
	bodyByte, _ := json.Marshal(car)

	r.req.Header.Add("Content-Type", "application/json")
	a := bytes.NewBuffer(bodyByte)
	rc := io.NopCloser(a)
	r.req.Body = rc

	return r
}

func (r *Request) UpdateCar(car *appsv1.Car) *Request {
	bodyByte, _ := json.Marshal(car)

	r.req.Header.Add("Content-Type", "application/json")
	a := bytes.NewBuffer(bodyByte)
	rc := io.NopCloser(a)
	r.req.Body = rc

	return r
}

func (r *Request) DeleteCar(name string) *Request {
	q := r.req.URL.Query()

	queryParam := map[string]string{
		"name": name,
	}
	for k, v := range queryParam {
		q.Add(k, v)
	}
	r.req.URL.RawQuery = q.Encode()

	return r
}

func (r *Request) ListCar() *Request {
	return r
}

func (r *Request) WatchCar() *Request {


	klog.Info("ws url:", r.url.String())
	// 创建ws连接
	c, _, err := websocket.DefaultDialer.Dial(r.url.String(), nil)
	// 赋值
	r.ws = c
	if err != nil {
		klog.Fatal("dial:", err)
	}
	//defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			// 读取对象
			_, message, err := r.ws.ReadMessage()

			if err != nil {
				klog.Errorf("read error:", err)
				return
			}

			// 放入chan中
			r.WChan <- message
		}
	}()

	return r
}



func (r *Request) GetAppleName(name string) *Request {

	q := r.req.URL.Query()

	queryParam := map[string]string{
		"name": name,
	}
	for k, v := range queryParam {
		q.Add(k, v)
	}
	r.req.URL.RawQuery = q.Encode()

	return r
}

func (r *Request) CreateApple(apple *v1.Apple) *Request {
	bodyByte, _ := json.Marshal(apple)

	r.req.Header.Add("Content-Type", "application/json")
	a := bytes.NewBuffer(bodyByte)
	rc := io.NopCloser(a)
	r.req.Body = rc

	return r
}

func (r *Request) UpdateApple(apple *v1.Apple) *Request {
	bodyByte, _ := json.Marshal(apple)

	r.req.Header.Add("Content-Type", "application/json")
	a := bytes.NewBuffer(bodyByte)
	rc := io.NopCloser(a)
	r.req.Body = rc

	return r
}

func (r *Request) DeleteApple(name string) *Request {
	q := r.req.URL.Query()

	queryParam := map[string]string{
		"name": name,
	}
	for k, v := range queryParam {
		q.Add(k, v)
	}
	r.req.URL.RawQuery = q.Encode()

	return r
}

func (r *Request) WsPath(p string) *Request {
	str := strings.Split(r.c.BasePath, "://")

	u := url.URL{Scheme: "ws", Host: str[1], Path: p}
	r.url = u

	return r
}

func (r *Request) WatchApple() *Request {


	klog.Info("ws url:", r.url.String())
	// 创建ws连接
	c, _, err := websocket.DefaultDialer.Dial(r.url.String(), nil)
	// 赋值
	r.ws = c
	if err != nil {
		klog.Fatal("dial:", err)
	}
	//defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			// 读取对象
			_, message, err := r.ws.ReadMessage()

			if err != nil {
				klog.Errorf("read error:", err)
				return
			}

			// 放入chan中
			r.WChan <- message
		}
	}()

	return r
}

func (r *Request) ListApple() *Request {
	return r
}



