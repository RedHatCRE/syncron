package logrus

import (
	"github.com/rhcre/syncron/pkg/cli"
	"github.com/sirupsen/logrus"
)

func Configure(cli cli.Input) error {
	steps := []func() error{
		// Set log level
		func() error {
			debug, err := cli.Cmd.Flags().GetBool("debug")

			if err != nil {
				return err
			}

			logrus.SetLevel(logrus.InfoLevel)

			if debug {
				logrus.SetLevel(logrus.TraceLevel)
			}

			return nil
		},
	}

	for _, step := range steps {
		err := step()

		if err != nil {
			return err
		}
	}

	return nil
}
