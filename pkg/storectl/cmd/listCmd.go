package cmds

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"log"
	"os"
	"practice_ctl/pkg/storectl/config"
	"practice_ctl/pkg/util/stores"
	"practice_ctl/pkg/util/stores/rest"
	"time"
)

// list操作 命令行 ex: storectl list apples 或 storectl list cars

var ListCmd = &cobra.Command{}

func ListCommand(configRes *config.StoreCtlConfig) *cobra.Command {
	// 配置文件

	cfg := &rest.Config{
		Host:    fmt.Sprintf("http://" + configRes.Server),
		Timeout: time.Second,
		Token:   configRes.Token,
	}
	clientSet := stores.NewForConfig(cfg)

	ListCmd = &cobra.Command{
		Use:          "list [flags]",
		Short:        "list [flags]",
		Example:      "storectl list apples",
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			// 区分输入命令 ex: storectl list apples 或是 storectl list car
			if args[0] == "apples" || args[0] == "apple" {
				err := ListApple(clientSet)
				if err != nil {
					return err
				}

			} else if args[0] == "cars" || args[0] == "car" {
				err := ListCar(clientSet)
				if err != nil {
					return err
				}
			}
			return nil
		},
	}

	return ListCmd

}

func ListApple(client *stores.ClientSet) error {

	appleList, err := client.CoreV1().Apple().List()
	if err != nil {
		log.Println(err)
		return err
	}
	// 表格化呈现
	table := tablewriter.NewWriter(os.Stdout)
	content := []string{"APPLE名称", "Price", "Place", "Color", "Size"}

	table.SetHeader(content)

	for _, apple := range appleList.Item {
		appleRow := []string{apple.Name, apple.Spec.Price, apple.Spec.Place, apple.Spec.Color, apple.Spec.Size}

		table.Append(appleRow)
	}
	// 去掉表格线
	//table = common.TableSet(table)

	table.Render()

	return nil
}

func ListCar(client *stores.ClientSet) error {

	carList, err := client.AppsV1().Car().List()
	if err != nil {
		log.Println(err)
		return err
	}
	// 表格化呈现
	table := tablewriter.NewWriter(os.Stdout)
	content := []string{"CAR名称", "Price", "Brand", "Color"}

	table.SetHeader(content)

	for _, car := range carList.Item {
		carRow := []string{car.Name, car.Spec.Price, car.Spec.Brand, car.Spec.Color}

		table.Append(carRow)
	}
	// 去掉表格线
	//table = common.TableSet(table)

	table.Render()

	return nil

}
