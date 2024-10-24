package cmd

import (
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of bmctl",
	Long:  `Print the version of bmctl`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println("bmctl version 0.1")
	},
}
