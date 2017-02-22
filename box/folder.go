package box

import (
	"encoding/json"
	"fmt"
)

func (c *client) GetFolder(id string) (*Folder, error) {
	req, err := c.Get("/folders/" + id)
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

func (c *client) GetFolderContents(id string) (*FolderContents, error) {
	req, err := c.Get("/folders/" + id + "/items" +
		"?fields=sequence_id,sha1,name,description,size," +
		"path_collection,created_at,modified_at,content_created_at," +
		"content_modified_at,created_by,modified_by,owned_by,parent," +
		"item_status,tags,has_collaborations,sync_status")
	if err != nil {
		return nil, err
	}

	var res Collection
	err = json.Unmarshal(req, &res)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(req))

	var files []File
	var folders []Folder
	for _, entry := range res.Entries {
		var entryType struct {
			Type string `json:"type"`
		}
		err := json.Unmarshal(entry, &entryType)
		if err != nil {
			return nil, err
		}

		switch entryType.Type {
		case FileType:
			var file File
			json.Unmarshal(entry, &file)
			files = append(files, file)
		case FolderType:
			var folder Folder
			json.Unmarshal(entry, &folder)
			folders = append(folders, folder)
		}
	}

	return &FolderContents{ID: id, Files: files, Folders: folders}, nil
}
