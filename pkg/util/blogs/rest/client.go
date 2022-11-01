package rest

import (
	"github.com/go-resty/resty/v2"
	"net/http"
)

type Interface interface {
	Get() *Request
	Post() *Request
}

var _ Interface = &RESTClient{}	// 查看是否实现此接口

// RESTClient 底层对象
type RESTClient struct {
	*resty.Client
}

// Get 方法
func (R *RESTClient) Get() *Request {
	return NewRequest(R).Verb(http.MethodGet)
}

// Post 方法
func (R *RESTClient) Post() *Request {
	return NewRequest(R).Verb(http.MethodPost)
}

// NewRESTClient 构建方法
func NewRESTClient(config *Config) *RESTClient {

	rc := resty.New()

	rc.SetBaseURL(config.Host)
	rc.SetTimeout(config.Timeout)

	return &RESTClient{Client: rc}
}


