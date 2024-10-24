package cmd

import (
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "bmctl",
		Short: "bmctl is a command line tool for interacting with BrandMeister API",
		Long:  `bmctl is a command line tool for interacting with BrandMeister API`,
	}
)

func Execute() error {
	return rootCmd.Execute()
}
