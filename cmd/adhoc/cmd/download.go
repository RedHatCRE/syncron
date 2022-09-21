package cmd

import (
	"github.com/rhcre/syncron/pkg/cli"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var download = &cobra.Command{
	Use:       "download",
	Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	ValidArgs: []string{cli.INSIGHTS_A, cli.SOSREPORTS_A, cli.ALL_A},
	RunE:      onRun,
}

func init() {
	download.Flags().Int(
		cli.DAYS_F,
		3,
		"Download data from the last X days",
	)

	download.Flags().Int(
		cli.MONTHS_F,
		0,
		"Download data from the last X months",
	)

	download.Flags().Int(
		cli.YEARS_F,
		0,
		"Download data from the last X years",
	)

	root.AddCommand(download)
}

func onRun(cmd *cobra.Command, args []string) error {
	logrus.Info("This is an info message.")
	logrus.Debug("This is a debug message.")
	logrus.Trace("This is a trace message, doubt they will happen much.")

	return nil
}
