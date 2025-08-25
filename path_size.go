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
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.Args().Len() == 0 {
				return fmt.Errorf("path argument is required")
			}

			path := cmd.Args().First()
			human := cmd.Bool("human")
            all := cmd.Bool("all")

			result, err := GetPathSize(path, false, human, all)
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
	if err != nil {
        return "", err
    }

    s := fmt.Sprintf("%s    %s", FormatSize(size, human), path)
    
    return s, err
}

func GetSize(path string, recursive, all bool) (int64, error) {
    var totalSize int64

    err := filepath.WalkDir(path, func(currentPath string, entry os.DirEntry, err error) error {
        if err != nil {
            return err
        }

        // Skip hidden files if  all == false
        if !all && entry.Name()[0] == '.' {
            if entry.IsDir() {
                return filepath.SkipDir
            }
            return nil
        }

        // If not in recursive mode, only files in the root of the directory are considered
        rel, err := filepath.Rel(path, currentPath)
		if err != nil {
    		return err
		}
		if !recursive && strings.Contains(rel, string(os.PathSeparator)) {
    		if entry.IsDir() {
        		return filepath.SkipDir
    		}
    		return nil
		}

        // Getting file size
        if !entry.IsDir() {
            info, err := entry.Info()
            if err != nil {
                return err
            }
            totalSize += info.Size()
        }
        return nil
    })

    return totalSize, err
}

func FormatSize(size int64, human bool) string {
    if !human {
        return fmt.Sprintf("%dB", size)
    }

    const unit = 1024
    if size < unit {
        return fmt.Sprintf("%dB", size)
    }

    div, exp := int64(unit), 0
    for n := size / unit; n >= unit; n /= unit {
        div *= unit
        exp++
    }

    units := []string{"KB", "MB", "GB", "TB", "PB", "EB"}
    return fmt.Sprintf("%.1f%s", float64(size)/float64(div), units[exp])
}


