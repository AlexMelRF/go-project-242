package path_size

import (
    "os"
    "github.com/urfave/cli/v3"
    "context"
)

func Run() {
    cmd := &cli.Command{
        Name:  "hexlet-path-size",
        Usage: "print size of a file or directory",
    }

    // Запускаем приложение
    err := cmd.Run(context.Background(), os.Args)
    if err != nil {
        panic(err)
    }
}