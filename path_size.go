package code

import (
    "os"
    "github.com/urfave/cli/v3"
    "context"
    "path/filepath"
    "fmt"
	"strings"
)

func Run() {
    // cmd := &cli.Command{
    //     Name:  "hexlet-path-size",
    //     Usage: "print size of a file or directory",
    // }

    // err := cmd.Run(context.Background(), os.Args)
    // if err != nil {
    //     // fmt.Println("hexlet-path-size: undefined operand \nTry hexlet-path-size --help for more information")
    // }
cmd := &cli.Command{
	Name:  "hexlet-path-size",
	Usage: "print size of a file or directory",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		if cmd.Args().Len() == 0 {
			return fmt.Errorf("path argument is required")
		}

		path := cmd.Args().First()
		result, err := GetPathSize(path, true, false, false)
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

func GetPathSize(path string, recursive, human, all bool) (string, error) {
    size, err := GetSize(path, recursive, all)
    s := ""
    if err == nil {
        s = fmt.Sprintf("%d    %s", size, path)
    }
    
    return s, err
}

func GetSize(path string, recursive, all bool) (int64, error) {
    var totalSize int64

    err := filepath.Walk(path, func(currentPath string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        // Skip hidden files if  all == false
        if !all && strings.HasPrefix(info.Name(), ".") {
            if info.IsDir() {
                return filepath.SkipDir
            }
            return nil
        }

        // If not in recursive mode, only files in the root of the directory are considered
        if !recursive {
            rel, _ := filepath.Rel(path, currentPath)
            if strings.Contains(rel, string(os.PathSeparator)) {
                if info.IsDir() {
                    return filepath.SkipDir
                }
                return nil
            }
        }

        if !info.IsDir() {
            totalSize += info.Size()
        }

        return nil
    })

    return totalSize, err
}

