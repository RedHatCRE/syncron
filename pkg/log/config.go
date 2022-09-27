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
func Configure(parser cli.CLIParser) error {
	// Define all actions needed to get Logrus ready to go
	setup := []func() error{
		// Set log level
		func() error {
			logrus.SetLevel(logrus.InfoLevel)

			if parser.GetDebug() {
				logrus.SetLevel(logrus.TraceLevel)
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
