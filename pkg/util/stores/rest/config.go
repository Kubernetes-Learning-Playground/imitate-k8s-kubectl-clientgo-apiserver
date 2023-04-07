package rest

import (
	"time"
)

// Config 配置
type Config struct {
	Host    string
	Timeout time.Duration
	Token   string
}


