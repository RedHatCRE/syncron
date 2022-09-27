package cli

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// CLIParser makes sense out of the data retrieved from a CLI framework.
//
// An implementation will ease access to the options the user has defined
// through the command line, such as the logging mode or the location of a
// configuration file.
type CLIParser interface {
	// GetDebug indicates whether the debug mode for Syncron has been
	// requested (true) or not (false).
	GetDebug() bool
}

// CobraParser takes care of extracting all interesting data from the
// structures provided by the Cobra CLI framework by implementing the
// cli.CLIParser interface.
type CobraParser struct {
	cmd  *cobra.Command
	args []string
}

// NewParserForCobra returns a new cli.CobraParser that will get to work with
// the Cobra structures given to this.
func NewParserForCobra(cmd *cobra.Command, args []string) *CobraParser {
	result := new(CobraParser)

	result.cmd = cmd
	result.args = args

	return result
}

// GetDebug implements cli.CLIParser.GetDebug by extracting the value of the
// 'debug' flag from Cobra's input data.
func (parser *CobraParser) GetDebug() bool {
	val, err := parser.cmd.Flags().GetBool(DEBUG_F)

	if err != nil {
		logrus.Errorf(
			"Unable to read value for argument: '%s'. Reason: '%s'.",
			DEBUG_F, err,
		)

		logrus.Warn(
			"Continuing without debug messages...",
		)

		return false
	}

	return val
}
