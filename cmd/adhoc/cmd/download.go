package cmd

import "github.com/spf13/cobra"

var download = &cobra.Command{
	Use:       "download",
	Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	ValidArgs: []string{"insights", "sosreports", "all"},
	RunE:      onRun,
}

func init() {
	download.Flags().Int("days", 3, "Download data from the last X days")
	download.Flags().Int("months", 0, "Download data from the last X months")
	download.Flags().Int("years", 0, "Download data from the last X years")

	root.AddCommand(download)
}

func onRun(cmd *cobra.Command, args []string) error {
	return nil
}
