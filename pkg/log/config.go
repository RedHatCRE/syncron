package log

import (
	"github.com/rhcre/syncron/pkg/cli"
	"github.com/sirupsen/logrus"
)

// Configure sets up Logrus so that it behaves as the user describes
// on the CLI.
//
// At this moment, Configure supports the following flags from the CLI:
//   - '-d'
//   - '--debug'
func Configure(in cli.Input) error {
	steps := []func() error{
		// Set log level
		func() error {
			debug, err := in.Cmd.Flags().GetBool(cli.DEBUG_F)

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
