package cmd

import "github.com/spf13/cobra"

var root = &cobra.Command{
	Use: "syncron",
}

func Execute() error {
	return root.Execute()
}

func init() {
	root.PersistentFlags().BoolP("debug", "d", false, "Turn on debug mode")
}
