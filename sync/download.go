package sync

import (
	"errors"
	"os"
	"path"

	"gitlab-beta.engr.illinois.edu/sp-box/boxsync/box"
)

const (
	syncRootName = "Box Sync"
)

func LocalSyncRoot() string {
	return path.Join(os.Getenv("HOME"), syncRootName)
}

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
