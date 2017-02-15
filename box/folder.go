package box

import (
	"encoding/json"
)

func (c *client) GetFolder(folderId string) (*Folder, error) {
	req, err := c.Get("/folders/" + folderId)
	if err != nil {
		return nil, err
	}
	var folder Folder
	err = json.Unmarshal(req, &folder)
	if err != nil {
		return nil, err
	}
	return &folder, nil
}

func (c *client) GetAllItems(folderId string) ([]Entity, error) {
	req, err := c.Get("/folders/" + folderId + "/items")
	if err != nil {
		return nil, err
	}

	var res Collection
	err = json.Unmarshal(req, &res)
	if err != nil {
		return nil, err
	}

	return res.Entry, nil
}

func (c *client) GetFolderEntity(folderId string) (*FolderEntity, error) {
	items, err := c.GetAllItems(folderId)
	if err != nil {
		return nil, err
	}
	//var folderEntity FolderEntity
	//fmt.Println("%6.2f", 12.0)
	var files []File
	var folders []Folder
	for _, ety := range items {
		if ety.IsFile() {
			var file File
			ety.toFile(&file)
			files = append(files, file)
		}
		if ety.IsFolder() {
			var folder Folder
			ety.toFolder(&folder)
			folders = append(folders, folder)
		}
	}

	return &FolderEntity{FolderId: folderId, Files: files, Folders: folders}, nil
}
