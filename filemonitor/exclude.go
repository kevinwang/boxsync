package filemonitor

import (
	"os"
	"path/filepath"
	"strings"
)

type Exclude struct {
	Patterns map[string]bool
	Files    map[string]bool
}

func (exclude *Exclude) IsMatch(filePath string) bool {
	return exclude.MatchPattern(filePath) || exclude.MatchFile(filePath)
}

func (exclude *Exclude) MatchPattern(filePath string) bool {
	baseName := filepath.Base(filePath)

	for pattern, _ := range exclude.Patterns {
		match, _ := filepath.Match(pattern, baseName)
		if match {
			return true
		}
	}

	return false
}

func (exclude *Exclude) MatchFile(filePath string) bool {
	filePath = filepath.Clean(filePath)
	fileInfo, _ := os.Stat(filePath)

	for excludePath, _ := range exclude.Files {
		/*
			checking if parent directory is excluded
			(1) Just checking by prefix now,
				(1.1) If the path is clean, it should be fine
			(2) maybe a traversal to root is better ?
		*/
		if strings.HasPrefix(filePath, excludePath) {
			return true
		}

		excludeInfo, _ := os.Stat(excludePath)
		//checking the same file
		if os.SameFile(fileInfo, excludeInfo) {
			return true
		}
	}

	return false
}
