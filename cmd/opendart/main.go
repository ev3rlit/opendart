package main

import (
	"fmt"
	"os"

	"github.com/awuzag/opendart/internal/cli"
)

func main() {
	cmd := cli.NewRootCommand(os.Stdout, os.Stderr, os.Getenv)
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
