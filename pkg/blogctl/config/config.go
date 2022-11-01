package config

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"practice_ctl/pkg/util/helpers"
)

// BlogCtlConfig
type BlogCtlConfig struct {
	Server string `yaml:"server"`
}

// LoadConfigFile 读取配置文件,模仿kubectl ，默认在~/.blog/config
func LoadConfigFile() *BlogCtlConfig {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln(err)
	}
	configFile := fmt.Sprintf("%s/.blog/config", home)
	if _, err := os.Stat(configFile); errors.Is(err, os.ErrNotExist) {
		log.Fatalln("配置文件没找到")
	}
	// 接配置文件
	cfg := &BlogCtlConfig{}
	err = yaml.Unmarshal(helpers.MustLoadFile(configFile), cfg)
	if err != nil {
		log.Fatalln(err)
	}
	return cfg

}
