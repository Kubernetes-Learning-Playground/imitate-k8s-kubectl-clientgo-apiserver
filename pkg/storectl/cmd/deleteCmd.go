package cmds

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	v1 "practice_ctl/pkg/apis/core/v1"
	"practice_ctl/pkg/storectl/config"
	"practice_ctl/pkg/util/stores"
	"practice_ctl/pkg/util/stores/rest"
	"strings"
	"time"
)

// TODO: delete操作

var DeleteCmd = &cobra.Command{}

func DeleteCommand(configRes *config.StoreCtlConfig) *cobra.Command {
	// 配置文件


	cfg := &rest.Config{
		Host:    fmt.Sprintf("http://" + configRes.Server),
		Timeout: time.Second,
		Token: configRes.Token,
	}
	clientSet := stores.NewForConfig(cfg)

	DeleteCmd = &cobra.Command{
		Use:          "delete [flags]",
		Short:        "delete [flags]",
		Example:      "storectl delete  apples",
		SilenceUsage: true,
		// args[0] 区分资源， args[1] json路径
		RunE: func(c *cobra.Command, args []string) error {
			// 模拟kubectl delete -f xxxx.yaml
			// 但是这个需要变成多一步
			// storectl delete cars --file xxxx.yaml

			if file != "" {
				// 区分yaml json

				s := strings.Split(file, ".")

				if args[0] == "apples" || args[0] == "apple" {
					err := DeleteApple(clientSet, file, s[len(s)-1])
					if err != nil {
						return err
					}

				} else if args[0] == "cars" || args[0] == "car" {
					err := DeleteCar(clientSet, file, s[len(s)-1])
					if err != nil {
						return err
					}
				}

			} else {
				// 这里就是 storectl delete car xxxxx
				if len(args) < 2 {
					return errors.New("input args least two")
				}
				if args[0] == "apples" || args[0] == "apple" {
					err := DeleteAppleWithName(clientSet, args[1])
					if err != nil {
						return err
					}

				} else if args[0] == "cars" || args[0] == "car" {
					err := DeleteCarWithName(clientSet, args[1])
					if err != nil {
						return err
					}
				}
			}


			return nil
		},
	}

	return DeleteCmd


}

// TODO: 目前删除只支持 yaml 或 json 两种方式，目前未支持输入name的删除方式

// DeleteApple 读取文件删除的方式
func DeleteApple(client *stores.ClientSet, path string, pathType string) error {

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("读取json文件失败", err)
		return err
	}
	a := &v1.Apple{}

	if pathType == "yaml" {
		err = yaml.Unmarshal(bytes, a)
	} else if pathType == "json"{
		err = json.Unmarshal(bytes, a)
	}

	if err != nil {
		fmt.Println("解析数据失败", err)
		return err
	}

	res, err := client.CoreV1().Apple().Get(a.Name)
	if res.Name == "" && err == nil {

		fmt.Printf("apple name:%v is notfound\n", a.Name)
		return nil
	}

	// 创建操作
	err = client.CoreV1().Apple().Delete(res.Name)
	if err != nil {
		fmt.Printf("apple name:%v delete error\n", res.Name)
		return nil
	}
	fmt.Printf("apple is delete\n")

	return nil

}


// DeleteCar 读取文件删除的方式
func DeleteCar(client *stores.ClientSet, path string, pathType string) error{
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("读取json文件失败", err)
		return err
	}
	a := &v1.Apple{}

	if pathType == "yaml" {
		err = yaml.Unmarshal(bytes, a)
	} else if pathType == "json"{
		err = json.Unmarshal(bytes, a)
	}

	if err != nil {
		fmt.Println("解析数据失败", err)
		return err
	}

	res, err := client.AppsV1().Car().Get(a.Name)
	if res.Name == "" && err == nil {

		fmt.Printf("car name:%v is notfound\n", a.Name)
		return nil
	}

	// 创建操作
	err = client.AppsV1().Car().Delete(res.Name)
	if err != nil {
		fmt.Printf("car name:%v delete error\n", res.Name)
		return nil
	}
	fmt.Printf("car is delete\n")

	return nil
}

// DeleteCarWithName 直接用name的方式删除
func DeleteCarWithName(client *stores.ClientSet, name string) error {
	res, err := client.AppsV1().Car().Get(name)
	if res.Name == "" && err == nil {

		fmt.Printf("car is notfound\n" )
		return nil
	}

	// 创建操作
	err = client.AppsV1().Car().Delete(res.Name)
	if err != nil {
		fmt.Printf("car name:%v delete error\n", res.Name)
		return nil
	}
	fmt.Printf("car is delete\n")


	return nil
}

// DeleteAppleWithName 直接用name的方式删除
func DeleteAppleWithName(client *stores.ClientSet, name string) error {
	res, err := client.CoreV1().Apple().Get(name)
	if res.Name == "" && err == nil {

		fmt.Printf("apple is notfound\n" )
		return nil
	}

	// 创建操作
	err = client.CoreV1().Apple().Delete(res.Name)
	if err != nil {
		fmt.Printf("apple name:%v delete error\n", res.Name)
		return nil
	}
	fmt.Printf("apple is delete\n")


	return nil
}