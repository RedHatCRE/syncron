package cmd

import "github.com/spf13/cobra"

var root = &cobra.Command{
	Use: "syncron",
}

func Execute() error {
	return root.Execute()
}

func init() {
}
