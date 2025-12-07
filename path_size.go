package code

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func GetPathSize(path string, recursive, human, all bool) (string, error) {
	size, err := getSize(path, recursive, all)

	if err != nil {
		return "", err
	}

	s := formatSize(size, human)

	return s, err
}

func getSize(path string, recursive, all bool) (int64, error) {
	var totalSize int64

	// get info about the root path
	rootInfo, err := os.Lstat(path)
	if err != nil {
		return 0, fmt.Errorf("cannot access '%s': %w", path, err)
	}

	// if it is a file, just return its size
	if !rootInfo.IsDir() {
		return rootInfo.Size(), nil
	}

	// if it's a dir - use filepath.WalkDir
	err = filepath.WalkDir(path, func(currentPath string, entry os.DirEntry, err error) error {
		if err != nil {
			// in case of a file access error, we display a warn and continue
			fmt.Fprintf(os.Stderr, "Warning: cannot access '%s': %v\n", currentPath, err)
			return nil
		}

		// get the relative path from the root dir
		rel, err := filepath.Rel(path, currentPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: cannot get relative path for '%s': %v\n", currentPath, err)
			return nil
		}

		// skip hidden files/dirs if all == false
		// check each path component for hidden status
		if !all {
			pathComponents := strings.Split(rel, string(os.PathSeparator))
			for _, component := range pathComponents {
				if component != "." && component != ".." && strings.HasPrefix(component, ".") {
					if entry.IsDir() {
						return filepath.SkipDir
					}
					return nil
				}
			}
		}

		// if not in recursive mode, skip subdirs and their contents
		if !recursive {
			// if it's a subdir (not the root), skip it
			if rel != "." && entry.IsDir() {
				return filepath.SkipDir
			}

			// if this is a file in a subdir, skip it
			if rel != "." && strings.Contains(rel, string(os.PathSeparator)) {
				return nil
			}
		}

		// get file info with os.Lstat
		info, err := os.Lstat(currentPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: cannot get file info for '%s': %v\n", currentPath, err)
			return nil
		}

		// only consider regular files (not dirs)
		if !info.IsDir() {
			totalSize += info.Size()
		}

		return nil
	})

	return totalSize, err
}

func formatSize(size int64, human bool) string {
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
