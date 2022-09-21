package cmd

import (
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

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
	log.Info("This is an info message.")
	log.Debug("This is a debug message.")
	log.Trace("This is a trace message, doubt they will happen much.")

	return nil
}
