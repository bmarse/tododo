package main

import (
	"context"
	"log"
	"os"

	"github.com/bmarse/tododo/internal/ui"
	"github.com/urfave/cli/v3"
)

var appVersion = "development" // This will be replaced by the build system with the actual version.

func main() {
	var todoFilename string
	app := &cli.Command{
		Name:      "tododo",
		Usage:     "The todo manager that should be extinct",
		UsageText: "tododo [options] FILE\n\nFILE is the file we will use to store and load todos.",
		Version:   appVersion,
		ArgsUsage: "FILE",
		Arguments: []cli.Argument{
			&cli.StringArg{
				Destination: &todoFilename,
				Value:       ".tododo.md",
			},
		},
		Action: func(context.Context, *cli.Command) error {
			return ui.Run(todoFilename)
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
