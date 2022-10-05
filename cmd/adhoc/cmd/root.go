package cmd

import (

	"github.com/rhcre/syncron/pkg/cli"
	"github.com/rhcre/syncron/pkg/log"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
		func() error {
			// Setting up file formatting
			viper.SetConfigFile("syncron.yaml")
			viper.SetConfigType("yaml")
			viper.AddConfigPath(".")
			// Reading from file
			err := viper.ReadInConfig()
			if err != nil {
				return err
			} else {
				logrus.Info("Your configuration file was read succesfully!")
				bucket_name := viper.Get("buckets.sosreports")
				logrus.Info("Reading from bucket: ", bucket_name)
			}
			return nil
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
