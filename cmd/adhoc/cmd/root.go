package cmd

import (
	"os"

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

func readFile() {
	// Setting up file formatting
	viper.SetConfigFile("syncron.yaml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	// Reading from file
	err := viper.ReadInConfig()
	if err != nil {
		logrus.Error("Config file not found.\n")
		logrus.Info("Make sure syncron.yaml is in the correct folder. \n")
		logrus.Info("Syncron looks for the config file on: \n")
		logrus.Info("	- Root of the repository \n")
		os.Exit(1)
	} else {
		logrus.Info("Your configuration file was read succesfully!")
	}
	private_key := viper.Get("credentials.aws_secret_access_key")
	public_key := viper.Get("credentials.aws_access_key_id")
	bucket_name := viper.Get("buckets.sosreports")
	logrus.Info("aws_secret_access_key: ", private_key)
	logrus.Info("aws_public_access_key: ", public_key)
	logrus.Info("Reading from bucket: ", bucket_name)
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
	readFile()
	for _, step := range setup {
		err := step()

		if err != nil {
			return err
		}
	}

	return nil
}
