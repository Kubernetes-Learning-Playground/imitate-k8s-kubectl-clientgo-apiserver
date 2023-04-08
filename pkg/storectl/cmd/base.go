package cmds

import (
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"log"
	"practice_ctl/pkg/storectl/config"
)

type CmdMetaData struct {
	Use     string
	Short   string
	Example string
}

var storeCmdMetaData *CmdMetaData

func init() {
	storeCmdMetaData = &CmdMetaData{
		Use:     "storectl [flags]",
		Short:   "模仿kubectl",
		Example: "storectl [flags]",
	}

}

// TODO: apple car命令行

// RunCmd
func RunCmd() {
	cmd := &cobra.Command{
		Use:          storeCmdMetaData.Use,
		Short:        storeCmdMetaData.Short,
		Example:      storeCmdMetaData.Example,
		SilenceUsage: true,
	}

	configRes := config.LoadConfigFile()

	// cmd操作
	listCmd := ListCommand(configRes)
	createCmd := CreateCommand(configRes)
	applyCmd := ApplyCommand(configRes)
	deleteCmd := DeleteCommand(configRes)
	describeCmd := DescribeCommand(configRes)
	//加入子命令

	deleteCmd.Flags().StringVar(&file, "file", "", "storectl delete cars --file xxxxx.yaml  ")

	cmd.AddCommand(versionCmd, listCmd, createCmd, applyCmd, deleteCmd, describeCmd)
	err := cmd.Execute()
	if err != nil {
		log.Fatalln(err)
	}
}

var file string

var cfgFlags *genericclioptions.ConfigFlags

func MergeFlags(cmds ...*cobra.Command) {
	cfgFlags = genericclioptions.NewConfigFlags(true)
	for _, cmd := range cmds {
		cfgFlags.AddFlags(cmd.Flags())
	}
}
