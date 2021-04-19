package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/nhatthm/n26cli/internal/app"
	"github.com/nhatthm/n26cli/internal/cli"
)

func main() {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	l := app.NewServiceLocator()

	if err := cli.NewApp(l, usr.HomeDir).Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Exit(0)
}
