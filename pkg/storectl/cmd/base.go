package cmds

import (
	"fmt"
	"github.com/spf13/cobra"
	"practice_ctl/pkg/storectl/config"
	"log"
)

// RunCmd
func RunCmd() {
	cmd := &cobra.Command{
		Use:          "storectl",
		Short:        "模仿kubectl",
		Example:      "storectl",
		SilenceUsage: true,
	}

	cfg := config.LoadConfigFile()
	fmt.Println(cfg)

	//加入子命令
	cmd.AddCommand(versionCmd)
	err := cmd.Execute()
	if err != nil {
		log.Fatalln(err)
	}
}
