package rest

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	v1 "practice_ctl/pkg/apis/core/v1"
	"strings"
)

// 模仿 k8s
type Request struct {
	c       *RESTClient
	req 	*http.Request

}

func NewRequest(c *RESTClient) *Request {
	req, _ := http.NewRequest("", "", nil)
	return &Request{c: c,
		req: req,

	} //默认是根
}

//最终执行 http请求
//func (r *Request) Do() (*resty.Response, error) {
//	return r.c.R().Execute(r.verb, r.path)
//}

// Do 需要对底层库进行封装，不要暴露
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

func (r *Request) ListApple() *Request {
	return r
}



