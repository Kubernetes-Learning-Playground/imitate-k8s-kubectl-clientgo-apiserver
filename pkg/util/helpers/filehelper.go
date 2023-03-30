package helpers

import (
	"io/ioutil"
	"os"
)

// MustLoadFile 如果读不到file，就panic
func MustLoadFile(path string) []byte {
	b, err := LoadFile(path)
	if err != nil {
		panic(err)
	}
	return b
}

// LoadFile 加载指定目录的文件, 全部取出内容
func LoadFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return b, err
}

func Must(err error) {
	if err != nil {
		panic(err)
	}
}
