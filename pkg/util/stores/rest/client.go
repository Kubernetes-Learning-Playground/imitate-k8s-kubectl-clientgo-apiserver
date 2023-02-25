package rest

import (
	"golang.org/x/net/websocket"
	"net/http"
)

type Interface interface {
	Get() *Request
	Post() *Request
	Put() *Request
	Delete() *Request
	Patch() *Request
}

var _ Interface = &RESTClient{}	// 查看是否实现此接口

// RESTClient 底层对象
type RESTClient struct {
	*http.Client
	Ws       *websocket.Conn
	BasePath string
}

// Get 方法
func (R *RESTClient) Get() *Request {
	return NewRequest(R).Verb(http.MethodGet)
}

// Post 方法
func (R *RESTClient) Post() *Request {
	return NewRequest(R).Verb(http.MethodPost)
}

// Put 方法
func (R *RESTClient) Put() *Request {
	return NewRequest(R).Verb(http.MethodPut)
}

// Patch 方法
func (R *RESTClient) Patch() *Request {
	return NewRequest(R).Verb(http.MethodPatch)
}

// Delete 方法
func (R *RESTClient) Delete() *Request {

	return NewRequest(R).Verb(http.MethodDelete)
}

// NewRESTClient 构建方法
func NewRESTClient(config *Config) *RESTClient {

	c := &http.Client{}
	c.Timeout = config.Timeout

	return &RESTClient{Client: c, BasePath: config.Host}
}


