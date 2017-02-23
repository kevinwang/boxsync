package box

import (
	"encoding/json"
)

func (c *client) GetFolder(id string) (*Folder, error) {
	body, err := c.Get("/folders/" + id)
	if err != nil {
		return nil, err
	}
	var folder Folder
	err = json.Unmarshal(body, &folder)
	if err != nil {
		return nil, err
	}
	return &folder, nil
}

func (c *client) GetFolderContents(id string) (*FolderContents, error) {
	body, err := c.Get("/folders/" + id + "/items" +
		"?fields=sequence_id,sha1,name,description,size," +
		"path_collection,created_at,modified_at,content_created_at," +
		"content_modified_at,created_by,modified_by,owned_by,parent," +
		"item_status,tags,has_collaborations,sync_status")
	if err != nil {
		return nil, err
	}

	var collection Collection
	err = json.Unmarshal(body, &collection)
	if err != nil {
		return nil, err
	}

	var files []File
	var folders []Folder
	for _, entry := range collection.Entries {
		var entryType struct {
			Type string `json:"type"`
		}
		err := json.Unmarshal(entry, &entryType)
		if err != nil {
			return nil, err
		}

		switch entryType.Type {
		case TypeFile:
			var file File
			json.Unmarshal(entry, &file)
			files = append(files, file)
		case TypeFolder:
			var folder Folder
			json.Unmarshal(entry, &folder)
			folders = append(folders, folder)
		}
	}

	return &FolderContents{ID: id, Files: files, Folders: folders}, nil
}
