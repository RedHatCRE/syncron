package cmd

import "github.com/spf13/cobra"

var download = &cobra.Command{
	Use:  "download",
	RunE: onRun,
}

func init() {
	root.AddCommand(download)
}

func onRun(cmd *cobra.Command, args []string) error {
	return nil
}
