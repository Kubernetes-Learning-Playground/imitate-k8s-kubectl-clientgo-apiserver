package cmds

import (
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	v12 "practice_ctl/pkg/apis/apps/v1"
	v1 "practice_ctl/pkg/apis/core/v1"
	"practice_ctl/pkg/storectl/config"
	"practice_ctl/pkg/util/stores"
	"practice_ctl/pkg/util/stores/rest"
	"time"
	"encoding/json"
)






var CreateCmd = &cobra.Command{}

func CreateCommand(configRes *config.StoreCtlConfig) *cobra.Command {
	// 配置文件


	cfg := &rest.Config{
		Host:    fmt.Sprintf("http://" + configRes.Server),
		Timeout: time.Second,
	}
	clientSet := stores.NewForConfig(cfg)

	CreateCmd = &cobra.Command{
		Use:          "create [flags]",
		Short:        "create [flags]]",
		Example:      "storectl create apples",
		SilenceUsage: true,
		// args[0] 区分资源， args[1] json路径
		RunE: func(c *cobra.Command, args []string) error {
			if args[0] == "apples" || args[0] == "apple" {
				err := CreateApple(clientSet, args[1])
				if err != nil {
					return err
				}

			} else if args[0] == "cars" || args[0] == "car" {
				err := CreateCar(clientSet, args[1])
				if err != nil {
					return err
				}
			}


			return nil
		},
	}

	return CreateCmd


}



func CreateApple(client *stores.ClientSet, path string) error {

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("读取json文件失败", err)
		return err
	}
	a := &v1.Apple{}
	err = json.Unmarshal(bytes, a)
	if err != nil {
		fmt.Println("解析数据失败", err)
		return err
	}

	// 创建操作
	apple, err := client.Core().Apple().Create(a)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("apple name:%v is created\n", apple.Name)

	return nil

}


func CreateCar(client *stores.ClientSet, path string) error{
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("读取json文件失败", err)
		return err
	}
	a := &v12.Car{}
	err = json.Unmarshal(bytes, a)
	if err != nil {
		fmt.Println("解析数据失败", err)
		return err
	}

	// 创建操作
	car, err := client.Apps().Car().Create(a)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("car name:%v is created\n", car.Name)

	return nil
}


