package main

import (
	"fmt"
	"os"

	"github.com/nhatthm/n26cli/internal/cli"
)

func main() {
	if err := cli.NewApp().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Exit(0)
}
