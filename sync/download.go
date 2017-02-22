package sync

import (
	"errors"
	"fmt"
	"os"
	"path"

	"gitlab-beta.engr.illinois.edu/sp-box/boxsync/box"
)

const (
	syncRootName = "Box Sync"
)

var (
	LocalSyncRoot = path.Join(os.Getenv("HOME"), syncRootName)
)

func GetSyncRootFolder(client box.Client) (*box.Folder, error) {
	contents, err := client.GetFolderContents("0")
	if err != nil {
		return nil, err
	}

	for _, folder := range contents.Folders {
		if folder.Name == "Box Sync" {
			return &folder, nil
		}
	}

	return nil, errors.New("Box Sync folder does not exist in Box root")
}

func DownloadAll(client box.Client, folderID, destPath string) error {
	fi, err := os.Stat(destPath)
	switch {
	case err != nil:
		return err
	case !fi.IsDir():
		return errors.New(destPath + " is not a directory")
	}

	fmt.Printf("Downloading folder %s to %s\n", folderID, destPath)
	contents, err := client.GetFolderContents(folderID)
	if err != nil {
		return err
	}

	for _, file := range contents.Files {
		filePath := path.Join(destPath, file.Name)

		if _, err := os.Stat(filePath); err == nil && SHA1(filePath) == file.SHA1 {
			fmt.Printf("Checksums match, skipping: %s\n", filePath)
			continue
		}

		fmt.Printf("Downloading file %s to %s\n", file.ID, filePath)
		err = client.DownloadFile(file.ID, filePath)
		if err != nil {
			return err
		}
	}

	for _, folder := range contents.Folders {
		folderPath := path.Join(destPath, folder.Name)

		if _, err := os.Stat(folderPath); os.IsNotExist(err) {
			fmt.Printf("Creating directory %s\n", folderPath)
			err := os.MkdirAll(folderPath, 0755)
			if err != nil {
				return err
			}
		}

		err = DownloadAll(client, folder.ID, folderPath)
		if err != nil {
			return err
		}
	}

	return nil
}
