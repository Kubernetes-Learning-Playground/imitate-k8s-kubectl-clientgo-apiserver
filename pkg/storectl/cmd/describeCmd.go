package cmds

import (
	"fmt"
	"github.com/spf13/cobra"
	"practice_ctl/pkg/storectl/config"
	"practice_ctl/pkg/util/stores"
	"practice_ctl/pkg/util/stores/rest"
	"time"
	"gopkg.in/yaml.v2"
)


var DescribeCmd = &cobra.Command{}

func DescribeCommand(configRes *config.StoreCtlConfig) *cobra.Command {
	// 配置文件


	cfg := &rest.Config{
		Host:    fmt.Sprintf("http://" + configRes.Server),
		Timeout: time.Second,
	}
	clientSet := stores.NewForConfig(cfg)

	DescribeCmd = &cobra.Command{
		Use:          "describe [flags]",
		Short:        "describe [flags]",
		Example:      "storectl describe  apples",
		SilenceUsage: true,

		RunE: func(c *cobra.Command, args []string) error {

			// 区分输入命令 ex: storectl describe apples xxxxx 或是 storectl describe car xxxxx
			if args[0] == "apples" || args[0] == "apple" {
				err := DescribeApple(clientSet, args[1])
				if err != nil {
					return err
				}

			} else if args[0] == "cars" || args[0] == "car" {
				err := DescribeCar(clientSet, args[1])
				if err != nil {
					return err
				}
			}


			return nil
		},
	}

	return DescribeCmd

}


func DescribeApple(client *stores.ClientSet, name string) error {
	res, err := client.CoreV1().Apple().Get(name)
	if res.Name == "" && err == nil {
		fmt.Printf("apple is notfound\n" )
		return nil
	}
	resByte, err := yaml.Marshal(res)
	if err != nil {
		fmt.Printf("apple name:%v describe error\n", res.Name)
		return nil
	}
	fmt.Printf(string(resByte))
	return nil
}

func DescribeCar(client *stores.ClientSet, name string) error {

	res, err := client.AppsV1().Car().Get(name)
	if res.Name == "" && err == nil {
		fmt.Printf("car is notfound\n" )
		return nil
	}
	resByte, err := yaml.Marshal(res)
	if err != nil {
		fmt.Printf("car name:%v describe error\n", res.Name)
		return nil
	}
	fmt.Printf(string(resByte))

	return nil
}





