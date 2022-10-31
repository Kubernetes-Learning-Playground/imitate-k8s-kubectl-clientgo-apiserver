package rest

import (
	"github.com/go-resty/resty/v2"
	"net/http"
)

type Interface interface {
	Get() *Request
	Post() *Request
}

var _ Interface = &RESTClient{}

// 本课程来自 程序员在囧途(www.jtthink.com) 咨询群：98514334
type RESTClient struct {
	*resty.Client
}

func (R *RESTClient) Get() *Request {
	return NewRequest(R).Verb(http.MethodGet)
}

func (R *RESTClient) Post() *Request {
	return NewRequest(R).Verb(http.MethodPost)
}

// 本课程来自 程序员在囧途(www.jtthink.com) 咨询群：98514334
func NewRESTClient(config *Config) *RESTClient {

	rc := resty.New()

	rc.SetBaseURL(config.Host)
	rc.SetTimeout(config.Timeout)

	return &RESTClient{Client: rc}
}

// 本课程来自 程序员在囧途(www.jtthink.com) 咨询群：98514334
