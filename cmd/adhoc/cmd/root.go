package cmd

import (
	"fmt"

	"github.com/redhatcre/syncron/pkg/cli"
	"github.com/redhatcre/syncron/pkg/log"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
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
			viper.SetConfigFile("config/syncron.yaml")
			viper.SetConfigType("yaml")
			viper.AddConfigPath(".")
			// Reading from file
			err := viper.ReadInConfig()
			if err != nil {
				return err
			} else {
				logrus.Info("Your configuration file was read succesfully")
				logrus.Info("Reading from bucket: ", viper.Get("bucket"))
			}
			return nil
		},
		func() error {
			// Initialize a session with AWS SDK
			// It will read from the file located at ~/.aws/credentials and syncron/config/syncron.yaml
			logrus.Info("Starting setup for AWS s3")
			sess, err := session.NewSession(&aws.Config{
				Region:   aws.String(viper.GetString("s3.region")),
				Endpoint: aws.String(viper.GetString("s3.endpoint")),
			},
			)
			if err != nil {
				fmt.Println("Error setting up AWS...", "Check config params", err)
			} else {
				logrus.Info("AWS setup succesful!")
			}
			logrus.Info("Accessing bucket...")
			svc := s3.New(sess)
			resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(viper.GetString("bucket"))})
			if err != nil {
				fmt.Println("Unable to list items in bucket", err, resp)
			} else {
				logrus.Info("Success getting into the bucket!")
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
