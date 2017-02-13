package box

import (
	"errors"
)

// IsFolder checks if the given entity is a folder
func (e *Entity) IsFolder() bool {
	if e.Type == "folder" {
		return true
	} else {
		return false
	}
}

// toFolder converts the given entity to a folder. Only attributes present in
// the entity are populated rest are untouched.
func (e *Entity) toFolder(f *Folder) error {
	if !e.IsFolder() {
		return errors.New("Entity is not a folder")
	}
	f.folderId = e.Id
	f.Name = e.Name
	f.ETag = e.ETag
	f.SequenceId = e.SequenceId
	return nil
}

// IsFile checks if the given entity is a file.
func (e *Entity) IsFile() bool {
	if e.Type == "file" {
		return true
	} else {
		return false
	}
}

// toFile converts the given entity to a file. Only attributes present in
// the entity are populated rest are untouched.
func (e *Entity) toFile(f *File) error {
	if !e.IsFile() {
		return errors.New("Entity is not a file")
	}
	f.fileId = e.Id
	f.Name = e.Name
	f.ETag = e.ETag
	f.SequenceId = e.SequenceId
	return nil
}
