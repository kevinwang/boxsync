package filemonitor

import (
	"os"
	"path/filepath"
)

func getSubFolders(filePath string) (dirs []string, err error) {
	err = filepath.Walk(filePath, func(newPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			return nil
		}

		dirs = append(dirs, newPath)
		return nil
	})

	return dirs, err
}
