package cli

import "github.com/spf13/cobra"

type Input struct {
	Cmd  *cobra.Command
	Args []string
}
