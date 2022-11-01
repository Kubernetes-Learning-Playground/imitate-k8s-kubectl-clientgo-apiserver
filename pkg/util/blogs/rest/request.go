package rest

import (
	"fmt"
	"net/http"
	"net/url"
)

// 模仿 k8s
type Request struct {
	c       *RESTClient
	path    string
	params  url.Values
	headers http.Header
	verb    string
}

func NewRequest(c *RESTClient) *Request {
	return &Request{c: c, path: "/"} //默认是根
}

//最终执行 http请求
//func (r *Request) Do() (*resty.Response, error) {
//	return r.c.R().Execute(r.verb, r.path)
//}

// Do 需要对底层库进行封装，不要暴露
func (r *Request) Do() Result {
	var ret Result
	rsp, err := r.c.R().Execute(r.verb, r.path)
	if err != nil {
		ret.err = err
	} else if rsp.IsError() {
		ret.err = fmt.Errorf("%v", rsp.Error())
	} else {
		ret.rsp = rsp
	}
	return ret
}

//  method 设置
func (r *Request) Verb(verb string) *Request {
	r.verb = verb
	return r
}
func (r *Request) Path(p string) *Request {
	r.path = p
	return r
}
