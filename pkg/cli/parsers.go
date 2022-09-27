package cli

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type CLIParser interface {
	GetDebug() bool
}

type CobraParser struct {
	cmd  *cobra.Command
	args []string
}

func NewParserForCobra(cmd *cobra.Command, args []string) *CobraParser {
	result := new(CobraParser)

	result.cmd = cmd
	result.args = args

	return result
}

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
