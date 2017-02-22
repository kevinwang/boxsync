package box

import (
	"encoding/json"
	"time"
)

const (
	FileType   = "file"
	FolderType = "folder"
)

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Login string `json:"login"`
}

type File struct {
	ID                string     `json:"id,omitempty"`                  // Box’s unique string identifying this file.
	SequenceID        string     `json:"sequence_id,omitempty"`         // A unique ID for use with the /events endpoint.
	ETag              string     `json:"etag,omitempty"`                // A unique string identifying the version of this file.
	SHA1              string     `json:"sha1,omitempty"`                // The sha1 hash of this file.
	Name              string     `json:"name,omitempty"`                // The name of this file.
	Description       string     `json:"description,omitempty"`         // The description of this file.
	Size              int        `json:"size,omitempty"`                // Size of this file in bytes.
	PathCollection    Collection `json:"path_collection,omitempty"`     // The path of folders to this item, starting at the root.
	CreatedAt         time.Time  `json:"created_at,omitempty"`          // When this file was created on Box’s servers.
	ModifiedAt        time.Time  `json:"modified_at,omitempty"`         // When this file was last updated on the Box servers.
	ContentCreatedAt  time.Time  `json:"content_created_at,omitempty"`  // When the content of this file was created.
	ContentModifiedAt time.Time  `json:"content_modified_at,omitempty"` // When the content of this file was last modified.
	CreatedBy         User       `json:"created_by,omitempty"`          // The user who first created file.
	ModifiedBy        User       `json:"modified_by,omitempty"`         // The user who last updated this file.
	OwnedBy           User       `json:"owned_by,omitempty"`            // The user who owns this file.
	Parent            *Folder    `json:"parent,omitempty"`              // The folder containing this file.
	ItemStatus        string     `json:"item_status,omitempty"`         // Whether this item is deleted or not.
	VersionNumber     string     `json:"version_number,omitempty"`      // The version of the file.
	CommentCount      int        `json:"comment_count,omitempty"`       // The number of comments on a file.
	Tags              []string   `json:"tags,omitempty"`                // All tags applied to this file.
	Extension         string     `json:"extension,omitempty"`           // Indicates the suffix, when available, on the file.
}

type Folder struct {
	ID                string     `json:"id,omitempty"`                  // The folder’s ID.
	SequenceID        string     `json:"sequence_id,omitempty"`         // A unique ID for use with the /events endpoint.
	ETag              string     `json:"etag,omitempty"`                // A unique string identifying the version of this folder.
	Name              string     `json:"name,omitempty"`                // The name of this folder.
	Description       string     `json:"description,omitempty"`         // The description of this folder.
	Size              int        `json:"size,omitempty"`                // Size of this file in bytes.
	PathCollection    Collection `json:"path_collection,omitempty"`     // The path of folders to this item, starting at the root.
	CreatedAt         time.Time  `json:"created_at,omitempty"`          // The time the folder was created.
	ModifiedAt        time.Time  `json:"modified_at,omitempty"`         // The time the folder or its contents were last modified.
	ContentCreatedAt  time.Time  `json:"content_created_at,omitempty"`  // The time the folder or its contents were originally created (according to the uploader).
	ContentModifiedAt time.Time  `json:"content_modified_at,omitempty"` // The time the folder or its contents were last modified (according to the uploader).
	CreatedBy         User       `json:"created_by,omitempty"`          // The user who created this folder.
	ModifiedBy        User       `json:"modified_by,omitempty"`         // The user who last modified this folder.
	OwnedBy           User       `json:"owned_by,omitempty"`            // The user who owns this file.
	Parent            *Folder    `json:"parent,omitempty"`              // The folder that contains this one.
	ItemStatus        string     `json:"item_status,omitempty"`         // Whether this item is deleted or not.
	Tags              []string   `json:"tags,omitempty"`                // All tags applied to this file.
	HasCollaborations bool       `json:"has_collaborations,omitempty"`  // Whether this folder has any collaborators.
	SyncStatus        string     `json:"sync_status,omitempty"`         // Whether this folder will be synced by the Box sync clients or not. Can be
}

type Collection struct {
	Count   int               `json:"total_count,omitempty"`
	Entries []json.RawMessage `json:"entries,omitempty"`
	Limit   int               `json:"limit,omitempty"`
	Offset  int               `json:"offset,omitempty"`
}

type FolderContents struct {
	ID      string `json:"id,omitempty"`
	Files   []File
	Folders []Folder
}
