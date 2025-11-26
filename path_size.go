package code

import (
    "os"
    "path/filepath"
    "fmt"
	"strings"
)

func GetPathSize(path string, recursive, human, all bool) (string, error) {
    size, err := GetSize(path, recursive, all)
	if err != nil {
        return "", err
    }

    s := FormatSize(size, human)
    
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


