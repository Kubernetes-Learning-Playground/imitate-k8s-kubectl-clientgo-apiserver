package rest

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
)

// Result 错误
type Result struct {
	rsp *resty.Response
	err error
}

// Into 把resp转出来
func (r Result) Into(v interface{}) error {
	if r.err != nil {
		return r.err
	}
	//TODO 这里简化了。 理论上应该做多格式的 解码器
	return json.Unmarshal(r.rsp.Body(), v)
}
