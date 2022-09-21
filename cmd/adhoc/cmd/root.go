package cmd

import (
	"github.com/rhcre/syncron/pkg/cli"
	"github.com/rhcre/syncron/pkg/log"
	"github.com/spf13/cobra"
)

var root = &cobra.Command{
	Use:               "syncron",
	PersistentPreRunE: onPersistentPreRun,
}

func Execute() error {
	return root.Execute()
}

func init() {
	root.PersistentFlags().BoolP(
		cli.DEBUG_F,
		cli.DEBUG_P,
		false,
		"Turn on debug mode",
	)
}

func onPersistentPreRun(cmd *cobra.Command, args []string) error {
	// Define actions taken to get Syncron ready before commands are executed
	setup := []func() error{
		// Set up logging
		func() error {
			return log.Configure(
				cli.Input{
					Cmd:  cmd,
					Args: args,
				},
			)
		},
	}

	for _, step := range setup {
		err := step()

		if err != nil {
			return err
		}
	}

	return nil
}
