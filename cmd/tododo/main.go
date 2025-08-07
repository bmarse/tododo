package main

import (
	"context"
	"log"
	"os"

	"github.com/bmarse/tododo/pkg/ui"
	"github.com/urfave/cli/v3"
)

var appVersion = "development" // This will be replaced by the build system with the actual version.

func main() {
	app := &cli.Command{
		Name:    "tododo",
		Usage:   "The todo manager that should be extinct",
		Version: appVersion,
		Action: func(context.Context, *cli.Command) error {
			return ui.Run()
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
