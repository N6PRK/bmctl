package main

import (
	"github.com/spf13/cobra"

	"github.com/n6prk/bmctl/cmd"
)

func main() {
	cobra.CheckErr(cmd.Execute())
}
