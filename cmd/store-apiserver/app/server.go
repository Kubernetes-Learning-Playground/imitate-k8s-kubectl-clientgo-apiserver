package app

import (
	"github.com/spf13/cobra"
)

func NewAPIServerCommand() *cobra.Command {
	s := NewServerRunOptions()
	cmd := &cobra.Command{
		Use: "store-apiserver",
		Short: "run http server",
		Long: ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			// set default options
			completedOptions, err := Complete(s)
			if err != nil {
				return err
			}

			// validate options
			if errs := completedOptions.Validate(); len(errs) != 0 {

				return errs[0]
			}
			stopC := make(chan struct{})
			return Run(completedOptions, stopC)
		},
	}
	// TODO: 这里可以让用户配置，使用 option -> config -> api-server
	//fs := cmd.Flags()
	//namedFlagSets := s.Flags()

	return cmd
}



