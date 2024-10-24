package main

import (
	"fmt"

	"github.com/n6prk/bmctl/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		panic(fmt.Sprintf("failed to execute command: %v", err))
	}
}
