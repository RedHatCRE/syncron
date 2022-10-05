package cmd

import (
	"github.com/redhatcre/syncron/pkg/cli"
	"github.com/redhatcre/syncron/pkg/log"
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
		cli.Debug,
		cli.D,
		false,
		"Turn on debug mode",
	)
}

func onPersistentPreRun(cmd *cobra.Command, args []string) error {
	// Define actions taken to get Syncron ready before commands are executed
	setup := []func() error{
		// Set up logging
		func() error {
			parser := cli.NewParserForCobra(cmd, args)

			return log.Configure(parser)
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
