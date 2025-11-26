package main

import (
		"code"
		"fmt"
		"github.com/urfave/cli/v3"
		"context"
		"os"
	)

func main() {
	cmd := &cli.Command {
		Name:  "hexlet-path-size",
		Usage: "print size of a file or directory",
		Flags: []cli.Flag {
			&cli.BoolFlag {
				Name:  "human",
				Aliases: []string{"H"},
				Usage: "human-readable sizes (auto-select unit)",
			},
            &cli.BoolFlag {
				Name:  "all",
				Aliases: []string{"a"},
				Usage: "include hidden files and directories",
			},
            &cli.BoolFlag {
				Name:  "recursive",
				Aliases: []string{"r"},
				Usage: "recursive size of directories",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.Args().Len() == 0 {
				return fmt.Errorf("path argument is required")
			}

			path := cmd.Args().First()
			human := cmd.Bool("human")
            all := cmd.Bool("all")
            recursive := cmd.Bool("recursive")

			result, err := code.GetPathSize(path, recursive, human, all)
			if err != nil {
				return err
			}

			fmt.Println(result)
			return nil
		},
	}
    if err := cmd.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
    }
}