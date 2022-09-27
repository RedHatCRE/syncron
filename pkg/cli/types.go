package cli

import "github.com/spf13/cobra"

// Input represents all data that the user provided through the CLI.
type Input struct {
	Cmd  *cobra.Command
	Args []string
}
