package main

import (
	"fmt"
	"os"

	"github.com/nhatthm/n26cli/internal/cli"
	"github.com/nhatthm/n26cli/internal/service"
)

func main() {
	l := &service.Locator{}

	if err := cli.NewApp(l).Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Exit(0)
}
