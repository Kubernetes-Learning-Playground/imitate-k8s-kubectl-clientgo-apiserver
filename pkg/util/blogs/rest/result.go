package rest

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
)

type Result struct {
	rsp *resty.Response
	err error
}

func (r Result) Into(v interface{}) error {
	if r.err != nil {
		return r.err
	}
	// 这块  是简化了 。  理论上 应该做多格式的 解码器
	return json.Unmarshal(r.rsp.Body(), v)
}
