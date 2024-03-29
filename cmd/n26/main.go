// Package main provides the main entry point for the n26 command line tool.
package main

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"

	"github.com/nhatthm/n26cli/internal/app"
	"github.com/nhatthm/n26cli/internal/cli"
)

func main() {
	homeDir, err := homedir.Dir()
	if err != nil {
		panic(err)
	}

	l := app.NewServiceLocator()

	if err := cli.NewApp(l, homeDir).Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Exit(0)
}
