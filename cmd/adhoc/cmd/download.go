package cmd

import (
	s3setup "github.com/redhatcre/syncron/pkg/bucketaws"
	"github.com/redhatcre/syncron/pkg/cli"
	"github.com/spf13/cobra"
)

var download = &cobra.Command{
	Use:       "download",
	Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	ValidArgs: []string{cli.Insights, cli.SOSReports, cli.All},
	RunE:      onRun,
}

func init() {
	download.Flags().Int(
		cli.Days,
		3,
		"Download data from the last X days",
	)

	download.Flags().Int(
		cli.Months,
		0,
		"Download data from the last X months",
	)

	download.Flags().Int(
		cli.Years,
		0,
		"Download data from the last X years",
	)

	root.AddCommand(download)
}

func onRun(cmd *cobra.Command, args []string) error {
	// Reading configuration file
	s3setup.ConfigRead()
	// Creating AWS session
	sess := s3setup.SetupSession()
	// Accessing bucket
	s3setup.AccessBucket(sess)
	return nil
}
